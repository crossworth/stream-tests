package main

import (
	"log"
	"net/http"
)

func main() {
	muxer := NewHLSMuxer("rtsp://170.93.143.139/rtplive/470011e600ef003a004ee33696235daa", "./hls/stream")
	muxer.Start()

	http.Handle("/", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		http.FileServer(http.Dir("./hls")).ServeHTTP(writer, request)
	}))

	log.Println("Starting WebServer at port 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
