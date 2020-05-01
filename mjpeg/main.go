package main

import (
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"
)

var currentImage []byte

func main() {
	muxer := NewMJPEGMuxer("rtsp://170.93.143.139/rtplive/470011e600ef003a004ee33696235daa", "./mjpeg/stream")
	muxer.Start()

	time.Sleep(1 * time.Second)

	go func() {
		for {
			imageData, err := ioutil.ReadFile("./mjpeg/stream/stream.jpg")
			if err != nil {
				log.Fatalln("error reading stream file", err)
				return
			}

			currentImage = imageData
			time.Sleep(1 * time.Second)
		}
	}()

	jpegHeader := make(textproto.MIMEHeader)
	jpegHeader.Add("Content-Type", "image/jpeg")

	http.Handle("/stream/stream.mjpg", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "multipart/x-mixed-replace;boundary=myboundary")
		writer.Header().Set("Cache-Control", "no-store")
		writer.Header().Set("Connection", "close")
		multipartWriter := multipart.NewWriter(writer)
		defer multipartWriter.Close()
		_ = multipartWriter.SetBoundary("myboundary")

		for {
			partWriter, _ := multipartWriter.CreatePart(jpegHeader)
			_, err := partWriter.Write(currentImage)
			if err != nil {
				log.Println("error writing to partWriter", err)
				return
			}

			time.Sleep(1 * time.Second)
		}
	}))

	http.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		http.FileServer(http.Dir("./mjpeg")).ServeHTTP(writer, request)
	}))

	log.Println("Starting WebServer at port 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
