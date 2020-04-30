package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const StreamMagicBytes = "jsmp"

var sizeRegex = regexp.MustCompile(`\d+x\d+`)

var clients map[string]*websocket.Conn
var width = uint16(0)
var height = uint16(0)

func main() {
	clients = make(map[string]*websocket.Conn)

	muxer := NewMpeg1Muxer("rtsp://170.93.143.139/rtplive/470011e600ef003a004ee33696235daa")

	muxer.OnData(func(data []byte) {
		for ip, c := range clients {
			log.Printf("Sending data to %s\n", ip)

			_ = c.SetWriteDeadline(time.Now().Add(time.Second * 5))
			err := c.WriteMessage(websocket.BinaryMessage, data)

			if err != nil {
				delete(clients, ip)
				return
			}

			log.Printf("Data sent to %s\n", ip)
		}
	})

	gettingInputData := false
	gettingOutputData := false
	var inputData []string

	muxer.OnError(func(data []byte) {
		input := string(data)
		if strings.Contains(input, "Input #") {
			gettingInputData = true
		}

		if strings.Contains(input, "Output #") {
			gettingInputData = false
			gettingOutputData = true
		}

		if gettingInputData {
			inputData = append(inputData, input)

			position := sizeRegex.FindStringIndex(input)

			if position != nil {
				sizeString := input[position[0]:position[1]]
				sizeParts := strings.Split(sizeString, "x")
				w, _ := strconv.Atoi(sizeParts[0])
				h, _ := strconv.Atoi(sizeParts[1])

				width = uint16(w)
				height = uint16(h)
				log.Printf("Size %dx%d\n", width, height)
			}
		}
	})

	log.Println("Starting Muxing")
	err := muxer.Start()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "index.html")
	})

	http.HandleFunc("/jsmpg.min.js", func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "jsmpg.min.js")
	})

	http.HandleFunc("/stream", handleWS)

	log.Println("Starting WebServer at port 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}

var upgrader = websocket.Upgrader{}

func handleWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("could not upgrade connection: %v\n", err)
		return
	}

	log.Printf("New WebSocket connection from: %s\n", r.RemoteAddr)

	var header bytes.Buffer
	header.WriteString(StreamMagicBytes)
	_ = binary.Write(&header, binary.BigEndian, width)
	_ = binary.Write(&header, binary.BigEndian, height)

	_ = c.SetWriteDeadline(time.Now().Add(time.Second * 5))
	err = c.WriteMessage(websocket.BinaryMessage, header.Bytes())
	if err != nil {
		log.Printf("could not write header: %v\n", err)
		_ = c.Close()
		return
	}

	log.Println("Header sent")
	clients[r.RemoteAddr] = c
}
