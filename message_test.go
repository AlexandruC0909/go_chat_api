package main

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestMessageEncode(t *testing.T) {
	id := uuid.New()
	message := &Message{
		Action:    SendMessageAction,
		Message:   "Hello, World!",
		Target:    &Room{Name: "Test Room"},
		Sender:    &Client{ID: id, Name: "Test User"},
		Timestamp: "2022-01-01T00:00:00Z",
	}

	jsonBytes := message.encode()
	var decodedMessage Message
	err := json.Unmarshal(jsonBytes, &decodedMessage)

	if err != nil {
		t.Fatalf("Failed to decode message: %v", err)
	}

	if decodedMessage.Action != message.Action {
		t.Errorf("Expected action %s, got %s", message.Action, decodedMessage.Action)
	}

	if decodedMessage.Message != message.Message {
		t.Errorf("Expected message %s, got %s", message.Message, decodedMessage.Message)
	}

	if decodedMessage.Target.Name != message.Target.Name {
		t.Errorf("Expected target room %s, got %s", message.Target.Name, decodedMessage.Target.Name)
	}

	if decodedMessage.Sender.ID != message.Sender.ID {
		t.Errorf("Expected sender ID %s, got %s", message.Sender.ID, decodedMessage.Sender.ID)
	}

	if decodedMessage.Sender.Name != message.Sender.Name {
		t.Errorf("Expected sender name %s, got %s", message.Sender.Name, decodedMessage.Sender.Name)
	}

	if decodedMessage.Timestamp != message.Timestamp {
		t.Errorf("Expected timestamp %s, got %s", message.Timestamp, decodedMessage.Timestamp)
	}
}
