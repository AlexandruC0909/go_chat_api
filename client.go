package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 10000
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn        *websocket.Conn
	wsServer    *WsServer
	send        chan []byte
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	rooms       map[*Room]bool
	RoomsIds    []uuid.UUID `json:"rooms"`
	isTyping    bool
	mu          sync.Mutex
	AvatarColor string `json:"avatarColor"`
}

func newClient(conn *websocket.Conn, wsServer *WsServer, name string) *Client {
	colors := []string{}
	prefixes := []string{
		"red",
		"cyan",
		"teal",
		"green",
		"orange",
		"deep-orange",
		"light-blue",
		"light-green",
		"lime",
		"amber",
		"deep-purple",
	}
	suffixes := []string{"9", "10"}

	for _, prefix := range prefixes {
		for _, suffix := range suffixes {
			color := fmt.Sprintf("%s-%s", prefix, suffix)
			colors = append(colors, color)
		}
	}
	return &Client{
		ID:          uuid.New(),
		Name:        name,
		conn:        conn,
		wsServer:    wsServer,
		send:        make(chan []byte, 256),
		rooms:       make(map[*Room]bool),
		RoomsIds:    make([]uuid.UUID, 0),
		AvatarColor: colors[rand.Intn(len(colors))],
	}

}

func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		client.handleNewMessage(jsonMessage)
	}

}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			client.mu.Lock()
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
			client.mu.Unlock()

		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) disconnect() {
	/* for room := range client.rooms {
		room.unregister <- client
	} */
	/* close(client.send)
	client.conn.Close()
	*/
	hasPrivateRoom := false
	for room := range client.rooms {
		if !room.Private {
			room.unregister <- client
		} else {
			hasPrivateRoom = true
		}
	}
	if !hasPrivateRoom {
		client.wsServer.unregister <- client
	}
	client.conn.Close()
}

func ServeWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {
	name, ok := r.URL.Query()["name"]
	if !ok || len(name[0]) < 1 {
		log.Println("Url Param 'name' is missing")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	id, ok := r.URL.Query()["id"]
	var client *Client
	if ok && len(id[0]) > 0 {
		client = wsServer.findClientByID(id[0])
		if client != nil {
			client.conn = conn
			go client.writePump()
			go client.readPump()
		} else {
			client = newClient(conn, wsServer, name[0])
			go client.writePump()
			go client.readPump()
		}
	} else {
		client = newClient(conn, wsServer, name[0])
		go client.writePump()
		go client.readPump()
	}

	roomListMsg := &RoomListMessage{
		Action:   "room-list",
		RoomList: wsServer.getAllRooms(client),
	}
	client.send <- roomListMsg.encode()

	message := &Message{
		Action: UserLoggedInAction,
		Sender: client,
	}
	client.send <- message.encode()

	if _, ok := wsServer.clients[client]; !ok {
		wsServer.register <- client
	}
	wsServer.listOnlineClients(client)
}

func (client *Client) handleNewMessage(jsonMessage []byte) {

	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}

	message.Sender = client
	currentTime := time.Now()
	currentHour, currentMinute, _ := currentTime.Clock()
	message.Timestamp = fmt.Sprintf("%d:%02d", currentHour, currentMinute)
	switch message.Action {
	case SendMessageAction:
		roomID := message.Target.GetId()
		if room := client.wsServer.findRoomByID(roomID); room != nil {
			room.Messages = append(room.Messages, message)
			room.broadcast <- &message
		}

	case JoinRoomAction:
		client.handleJoinRoomMessage(message)

	case LeaveRoomAction:
		client.handleLeaveRoomMessage(message)

	case JoinRoomPrivateAction:
		client.handleJoinRoomPrivateMessageSimple(message)

	case TypingAction:
		client.SetTyping(message)

	case DeleteRoomAction:
		client.handleDeleteRoomAcion(message)

	}

}

func (client *Client) handleDeleteRoomAcion(message Message) {
	room := client.wsServer.findRoomByID(message.Target.GetId())
	if room == nil {
		return
	}
	client.wsServer.deleteRoom(room)
	roomListMsg := &RoomListMessage{
		Action:   "room-list",
		RoomList: client.wsServer.getAllRooms(client),
	}
	for otherClients := range client.wsServer.clients {
		otherClients.send <- roomListMsg.encode()
	}
}

func (client *Client) handleJoinRoomMessage(message Message) {
	roomName := message.Message

	client.joinRoom(roomName, message.Sender, false)

}

func (client *Client) handleLeaveRoomMessage(message Message) {
	room := client.wsServer.findRoomByID(message.Message)
	if room == nil {
		return
	}

	delete(client.rooms, room)

	room.unregister <- client

}
func (client *Client) handleJoinRoomPrivateMessageSimple(message Message) {

	target := client.wsServer.findClientByID(message.Message)

	if target == nil {
		return
	}

	roomName := message.Message + client.ID.String()

	client.joinRoom(roomName, target, true)
	target.joinRoom(roomName, client, true)

}

func (client *Client) joinRoom(roomName string, sender *Client, private bool) {
	room := client.wsServer.findRoomByName(roomName)
	if room == nil {
		room = client.wsServer.createRoom(roomName, private, sender)
		for otherClients := range client.wsServer.clients {
			roomListMsg := &RoomListMessage{
				Action:   "room-list",
				RoomList: client.wsServer.getAllRooms(client),
			}
			otherClients.send <- roomListMsg.encode()
		}

	}

	if sender == nil && room.Private && !room.clients[client] {
		return
	}

	room.registerClientInRoom(sender)

	if !client.isInRoom(room) {

		client.rooms[room] = true

	}
	room.register <- client

	client.getRoomClients(room)
	client.notifyRoomJoined(room, sender)

}

func (client *Client) isInRoom(room *Room) bool {
	if _, ok := client.rooms[room]; ok {
		return true
	}

	return false
}

func (client *Client) notifyRoomJoined(room *Room, sender *Client) {
	message := Message{
		Action: RoomJoinedAction,
		Target: room,
		Sender: sender,
	}

	client.send <- message.encode()
}

func (client *Client) GetName() string {
	return client.Name
}

func (client *Client) SetTyping(message Message) {
	client.isTyping = message.Message == "true"

	message.Sender = client
	for room := range client.rooms {
		room.broadcast <- &message
	}
}

func contains(slice []uuid.UUID, item uuid.UUID) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func (client *Client) getRoomClients(room *Room) {

	for otherClient := range room.clients {

		roomListMsg := &RoomClientsListMessage{
			Action:          "room-clients-list",
			RoomClientsList: room.Clients,
		}
		otherClient.send <- roomListMsg.encode()

	}
}
