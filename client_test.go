package main

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestRegisterClientSuccessfully(t *testing.T) {
	server := NewWebsocketServer()
	client := &Client{ID: uuid.New()}

	go server.Run()

	server.register <- client

	time.Sleep(100 * time.Millisecond)

	if _, exists := server.clients[client]; !exists {
		t.Errorf("Expected client to be registered, but it was not")
	}
}
