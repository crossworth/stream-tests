package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var clients map[string]*websocket.Conn

func main() {
	clients = make(map[string]*websocket.Conn)

	muxer := NewFMP4Muxer("rtsp://170.93.143.139/rtplive/470011e600ef003a004ee33696235daa")

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

	muxer.OnError(func(data []byte) {
		input := string(data)
		fmt.Println(input)
	})

	log.Println("Starting Muxing")
	err := muxer.Start()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		http.FileServer(http.Dir("./fmp4")).ServeHTTP(writer, request)
	}))

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
	clients[r.RemoteAddr] = c
}
