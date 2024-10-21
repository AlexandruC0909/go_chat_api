package main

import (
	"encoding/json"
	"log"
)

const SendMessageAction = "send-message"
const SendAudioMessageAction = "send-audio-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"
const UserJoinedAction = "user-join"
const UserLeftAction = "user-left"
const JoinRoomPrivateAction = "join-room-private"
const RoomJoinedAction = "room-joined"
const TypingAction = "typing-action"
const UserLoggedInAction = "user-logged-in"
const DeleteRoomAction = "delete-room"

type Message struct {
	Action    string  `json:"action"`
	Message   string  `json:"message"`
	Target    *Room   `json:"target"`
	Sender    *Client `json:"sender"`
	Timestamp string  `json:"timestamp"`
	AudioData []byte  `json:"audioData"`
}
type RoomListMessage struct {
	Action   string  `json:"action"`
	RoomList []*Room `json:"rooms"`
}
type RoomClientsListMessage struct {
	Action          string    `json:"action"`
	RoomClientsList []*Client `json:"clients"`
}
type ClientsListMessage struct {
	Action      string    `json:"action"`
	ClientsList []*Client `json:"clients"`
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}

func (roomListMessage *RoomListMessage) encode() []byte {
	json, err := json.Marshal(roomListMessage)
	if err != nil {
		log.Println(err)
	}

	return json
}

func (roomListMessage *RoomClientsListMessage) encode() []byte {
	json, err := json.Marshal(roomListMessage)
	if err != nil {
		log.Println(err)
	}

	return json
}

func (clientsListMessage *ClientsListMessage) encode() []byte {
	json, err := json.Marshal(clientsListMessage)
	if err != nil {
		log.Println(err)
	}

	return json
}
