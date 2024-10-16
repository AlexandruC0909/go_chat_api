package main

import (
	"sync"
	"testing"

	"github.com/google/uuid"
)

func TestNewWebsocketServer(t *testing.T) {
	server := NewWebsocketServer()
	if server == nil {
		t.Error("Expected a new WsServer instance, got nil")
		return
	}
	if server.clients == nil {
		t.Error("Expected clients map to be initialized")
	}
	if server.register == nil {
		t.Error("Expected register channel to be initialized")
	}
	if server.unregister == nil {
		t.Error("Expected unregister channel to be initialized")
	}
	if server.broadcast == nil {
		t.Error("Expected broadcast channel to be initialized")
	}
	if server.rooms == nil {
		t.Error("Expected rooms map to be initialized")
	}
}

func TestRegisterClient(t *testing.T) {
	server := NewWebsocketServer()
	client := &Client{}
	server.registerClient(client)
	if _, ok := server.clients[client]; !ok {
		t.Error("Expected client to be registered")
	}
}

func TestUnregisterClient(t *testing.T) {
	server := NewWebsocketServer()
	client := &Client{}
	server.clients[client] = true
	server.unregisterClient(client)
	if _, ok := server.clients[client]; ok {
		t.Error("Expected client to be unregistered")
	}
}

func TestFindRoomByName(t *testing.T) {
	server := NewWebsocketServer()
	room := &Room{}
	room.Name = "test"
	server.rooms[room] = true
	foundRoom := server.findRoomByName("test")
	if foundRoom != room {
		t.Error("Expected to find the room")
	}
}

func TestFindRoomByID(t *testing.T) {
	server := NewWebsocketServer()
	room := &Room{}
	id := uuid.New()
	room.ID = id
	server.rooms[room] = true
	foundRoom := server.findRoomByID(id.String())
	if foundRoom != room {
		t.Error("Expected to find the room")
	}
}

func TestCreateRoom(t *testing.T) {
	server := NewWebsocketServer()
	room := server.createRoom("test", false, newClient(nil, nil, "test"))
	if room == nil {
		t.Error("Expected a new Room instance, got nil")
	}
	if _, ok := server.rooms[room]; !ok {
		t.Error("Expected room to be added to the server's rooms")
	}
}

func TestFindClientByID(t *testing.T) {
	server := NewWebsocketServer()
	client := &Client{}
	id := uuid.New()
	client.ID = id
	server.clients[client] = true
	foundClient := server.findClientByID(id.String())
	if foundClient != client {
		t.Error("Expected to find the client")
	}
}

func TestGetAllRooms(t *testing.T) {
	server := NewWebsocketServer()
	room1 := &Room{}
	room2 := &Room{}
	server.rooms[room1] = true
	server.rooms[room2] = true
	rooms := server.getAllRooms(newClient(nil, nil, "test"))
	if len(rooms) != 2 {
		t.Errorf("Expected 2 rooms, got %d", len(rooms))
	}
	for _, r := range rooms {
		if r != room1 && r != room2 {
			t.Error("Expected rooms to be equal")
		}
	}
}

func TestGetAllRoomsConcurrent(t *testing.T) {
	server := NewWebsocketServer()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		server.getAllRooms(newClient(nil, nil, "test"))
	}()
	go func() {
		defer wg.Done()
		server.getAllRooms(newClient(nil, nil, "test"))
	}()
	wg.Wait()
}
