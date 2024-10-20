package main

import (
	"github.com/google/uuid"
)

type Room struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Clients    []*Client `json:"clients"`
	Owner      *Client   `json:"owner"`
	Messages   []Message `json:"messages"`
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	Private    bool `json:"private"`
}

func NewRoom(name string, private bool, owner *Client) *Room {
	return &Room{
		ID:         uuid.New(),
		Name:       name,
		clients:    make(map[*Client]bool),
		Owner:      owner,
		Messages:   make([]Message, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
		Private:    private,
		Clients:    make([]*Client, 0),
	}
}

func (room *Room) RunRoom() {
	for {
		select {

		case client := <-room.register:
			room.registerClientInRoom(client)

		case client := <-room.unregister:
			room.unregisterClientInRoom(client)

		case message := <-room.broadcast:
			room.broadcastToClientsInRoom(message.encode())
		}

	}
}

func (room *Room) registerClientInRoom(client *Client) {
	if _, ok := room.clients[client]; !ok {
		room.clients[client] = true
		room.Clients = append(room.Clients, client)
	}
}

func (room *Room) unregisterClientInRoom(client *Client) {

	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)

		for i, existingClient := range room.Clients {
			if existingClient.ID == client.ID {
				room.Clients = append(room.Clients[:i], room.Clients[i+1:]...)
				break
			}
		}
	}

	client.getRoomClients(room)

}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.send <- message
	}
}

func (room *Room) GetId() string {
	return room.ID.String()
}

func (room *Room) GetName() string {
	return room.Name
}

func (room *Room) hasClient(client *Client) bool {
	_, ok := room.clients[client]
	return ok
}
