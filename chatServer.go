package main

import (
	"log"
	"sort"
	"sync"
)

type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	rooms      map[*Room]bool
	mutex      sync.Mutex
}

func NewWebsocketServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		rooms:      make(map[*Room]bool),
	}
}

func (server *WsServer) Run() {
	for {
		select {
		case client := <-server.register:
			server.registerClient(client)
			log.Printf("Client registered: %v", client)

		case client := <-server.unregister:
			server.unregisterClient(client)
			log.Printf("Client unregistered: %v", client)

		case message := <-server.broadcast:
			server.broadcastToClients(message)
			log.Printf("Broadcast message: %v", message)

		}
	}
}

func (server *WsServer) registerClient(client *Client) {
	server.clients[client] = true
	server.listOnlineClients()
}

func (server *WsServer) unregisterClient(client *Client) {

	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
	}
	server.listOnlineClients()

}

func (server *WsServer) listOnlineClients() {
	var clientList []*Client

	for otherClient := range server.clients {
		clientList = append(clientList, otherClient)
	}

	for otherClient := range server.clients {
		roomListMsg := &ClientsListMessage{
			Action:      UserJoinedAction,
			ClientsList: clientList,
		}
		otherClient.send <- roomListMsg.encode()
	}
}

func (server *WsServer) broadcastToClients(message []byte) {
	for client := range server.clients {
		client.send <- message
	}
}

func (server *WsServer) findRoomByName(name string) *Room {
	var foundRoom *Room
	for room := range server.rooms {
		if room.GetName() == name {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (server *WsServer) findRoomByID(ID string) *Room {
	var foundRoom *Room
	for room := range server.rooms {
		if room.GetId() == ID {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (server *WsServer) createRoom(name string, private bool, owner *Client) *Room {
	room := NewRoom(name, private, owner)
	go room.RunRoom()
	server.rooms[room] = true

	return room
}

func (server *WsServer) deleteRoom(room *Room) {
	delete(server.rooms, room)
}

func (server *WsServer) findClientByID(ID string) *Client {
	var foundClient *Client
	for client := range server.clients {
		if client.ID.String() == ID {
			foundClient = client
			break
		}
	}

	return foundClient
}

func (server *WsServer) getAllRooms(client *Client) []*Room {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	rooms := make([]*Room, 0, len(server.rooms))

	for room := range server.rooms {
		if !room.Private || (room.Private && room.clients[client]) {
			rooms = append(rooms, room)
		}
	}

	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].Name < rooms[j].Name
	})

	return rooms
}
