package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

const welcomeMessage = "%s joined the room"
const goodbyeMessage = "%s left the room"

type Room struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Clients    []Client  `json:"clients"`
	Messages   []Message `json:"messages"`
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	Private    bool `json:"private"`
}

type RoomListMessage struct {
	Action   string  `json:"action"`
	RoomList []*Room `json:"rooms"`
}

func NewRoom(name string, private bool) *Room {
	return &Room{
		ID:         uuid.New(),
		Name:       name,
		clients:    make(map[*Client]bool),
		Messages:   make([]Message, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
		Private:    private,
		Clients:    make([]Client, 0),
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
	room.notifyClientJoined(client)

	if _, ok := room.clients[client]; !ok {
		room.clients[client] = true
		room.Clients = append(room.Clients, *client)
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

	//room.notifyClientLeft(client)
}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.send <- message
	}
}

func (room *Room) notifyClientJoined(client *Client) {
	if !room.Private {
		currentTime := time.Now()
		currentHour, currentMinute, _ := currentTime.Clock()
		message := &Message{
			Action:    SendMessageAction,
			Target:    room,
			Sender:    client,
			Message:   fmt.Sprintf(welcomeMessage, client.GetName()),
			Timestamp: fmt.Sprintf("%d:%02d", currentHour, currentMinute),
		}

		room.broadcastToClientsInRoom(message.encode())
	}

}
func (room *Room) notifyClientLeft(client *Client) {
	currentTime := time.Now()
	currentHour, currentMinute, _ := currentTime.Clock()
	message := &Message{
		Action:    SendMessageAction,
		Target:    room,
		Sender:    client,
		Message:   fmt.Sprintf(goodbyeMessage, client.GetName()),
		Timestamp: fmt.Sprintf("%d:%02d", currentHour, currentMinute),
	}

	room.broadcastToClientsInRoom(message.encode())
}

func (room *Room) GetId() string {
	return room.ID.String()
}

func (room *Room) GetName() string {
	return room.Name
}

func (msg *RoomListMessage) encode() []byte {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error on encoding RoomListMessage: %s", err)
	}
	return data
}

func (room *Room) hasClient(client *Client) bool {
	_, ok := room.clients[client]
	return ok
}
