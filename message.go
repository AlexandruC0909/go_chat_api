package main

import (
	"encoding/json"
	"log"
)

const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"
const UserJoinedAction = "user-join"
const UserLeftAction = "user-left"
const JoinRoomPrivateAction = "join-room-private"
const RoomJoinedAction = "room-joined"
const TypingAction = "typing-action"

type Message struct {
	Action    string  `json:"action"`
	Message   string  `json:"message"`
	Target    *Room   `json:"target"`
	Sender    *Client `json:"sender"`
	Timestamp string  `json:"timestamp"`
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}
