package main

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewRoom(t *testing.T) {
	room := NewRoom("TestRoom", false)

	assert.NotNil(t, room)
	assert.Equal(t, "TestRoom", room.Name)
	assert.False(t, room.Private)
	assert.NotNil(t, room.clients)
	assert.NotNil(t, room.register)
	assert.NotNil(t, room.unregister)
	assert.NotNil(t, room.broadcast)
}

func TestRoom_GetId(t *testing.T) {
	room := NewRoom("TestRoom", false)

	assert.NotEmpty(t, room.GetId())
}

func TestRoom_GetName(t *testing.T) {
	room := NewRoom("TestRoom", false)

	assert.Equal(t, "TestRoom", room.GetName())
}

func TestRoomListMessage_encode(t *testing.T) {
	msg := &RoomListMessage{
		Action:   "testAction",
		RoomList: []*Room{NewRoom("TestRoom", false)},
	}

	data := msg.encode()
	assert.NotNil(t, data)

	var decodedMsg RoomListMessage
	err := json.Unmarshal(data, &decodedMsg)
	assert.NoError(t, err)
	assert.Equal(t, msg.Action, decodedMsg.Action)
	assert.Equal(t, len(msg.RoomList), len(decodedMsg.RoomList))
}

func TestRoom_registerClientInRoom(t *testing.T) {
	room := NewRoom("TestRoom", false)
	client := &Client{ID: uuid.New()}

	room.registerClientInRoom(client)

	assert.Contains(t, room.clients, client)
}

func TestRoom_unregisterClientInRoom(t *testing.T) {
	room := NewRoom("TestRoom", false)
	client := &Client{ID: uuid.New()}

	room.registerClientInRoom(client)
	room.unregisterClientInRoom(client)

	assert.NotContains(t, room.clients, client)
}
