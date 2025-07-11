package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var addr = flag.String("addr", ":8085", "http server address")

func main() {
	flag.Parse()

	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	wsServer := NewWebsocketServer()
	go func() {
		log.Println("Starting WebSocket server...")
		wsServer.Run()
	}()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received WebSocket connection request")
		ServeWs(wsServer, w, r)
	})

	log.Println("Starting HTTP server on", *addr)
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("HTTP server failed to start:", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Shutting down server...")
		os.Exit(0)
	}()
}
