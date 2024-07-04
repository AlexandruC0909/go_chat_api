package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8085", "http server address")

func main() {
	flag.Parse()

	wsServer := NewWebsocketServer()
	go wsServer.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(wsServer, w, r)
	})

	log.Fatal(http.ListenAndServe(*addr, nil))
}
