package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
)

const welcomeMessage = "%s joined the room"

type Room struct {
	ID         uuid.UUID   `json:"id"`
	Name       string      `json:"name"`
	ClientIDs  []uuid.UUID `json:"client_ids"`
	Messages   []Message   `json:"messages"`
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

// NewRoom creates a new Room
func NewRoom(name string, private bool) *Room {
	return &Room{
		ID:         uuid.New(),
		Name:       name,
		clients:    make(map[*Client]bool),
		Messages:   make([]Message, 0),
		ClientIDs:  make([]uuid.UUID, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
		Private:    private,
	}
}

// RunRoom runs our room, accepting various requests
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
	if !room.Private {
		room.notifyClientJoined(client)
	}
	room.clients[client] = true
	room.ClientIDs = append(room.ClientIDs, client.ID)
}

func (room *Room) unregisterClientInRoom(client *Client) {
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)

		for i, id := range room.ClientIDs {
			if id == client.ID {
				room.ClientIDs = append(room.ClientIDs[:i], room.ClientIDs[i+1:]...)
				break
			}
		}
	}
}

func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.send <- message
	}
}

func (room *Room) notifyClientJoined(client *Client) {
	message := &Message{
		Action:  SendMessageAction,
		Target:  room,
		Sender:  client,
		Message: fmt.Sprintf(welcomeMessage, client.GetName()),
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
