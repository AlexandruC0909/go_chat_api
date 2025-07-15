# Go-Chat

A simple chat application built with Go and WebSockets.

## Table of Contents

- [Go-Chat](#go-chat)
  - [Table of Contents](#table-of-contents)
  - [Introduction](#introduction)
  - [Features](#features)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
  - [Usage](#usage)
    - [Running the server](#running-the-server)
    - [Connecting a client](#connecting-a-client)
  - [API](#api)
    - [WebSocket Connection](#websocket-connection)
    - [Message Actions](#message-actions)
      - [send-message](#send-message)
      - [send-audio-message](#send-audio-message)
      - [join-room](#join-room)
      - [leave-room](#leave-room)
      - [user-join](#user-join)
      - [user-left](#user-left)
      - [join-room-private](#join-room-private)
      - [room-joined](#room-joined)
      - [typing-action](#typing-action)
      - [user-logged-in](#user-logged-in)
      - [delete-room](#delete-room)
  - [Project Structure](#project-structure)
  - [Contributing](#contributing)
  - [License](#license)

## Introduction

Go-Chat is a real-time chat application that allows users to communicate with each other through public and private rooms. It uses Go on the backend and WebSockets for real-time communication.

## Features

- Public and private chat rooms
- Real-time messaging
- Typing indicators
- User online status
- Audio messaging

## Getting Started

### Prerequisites

- Go 1.15 or higher

### Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/AlexandruC0909/go_chat_api.git
   ```
2. Navigate to the project directory:
   ```sh
   cd go_chat_api
   ```
3. Install dependencies:
   ```sh
   go mod tidy
   ```

## Usage

### Running the server

To run the chat server, execute the following command:

```sh
go run .
```

By default, the server will start on port 8085. You can change the port by using the `-addr` flag:

```sh
go run . -addr :8080
```

### Connecting a client

To connect a client to the server, you need to establish a WebSocket connection to the `/ws` endpoint. You must provide a `name` query parameter for the client's name.

Example: `ws://localhost:8085/ws?name=JohnDoe`

## API

### WebSocket Connection

- **Endpoint**: `/ws`
- **Query Parameters**:
  - `name` (string, required): The name of the client.
  - `id` (string, optional): The ID of an existing client to reconnect.

### Message Actions

The `action` field in the JSON message determines the type of action to be performed.

#### send-message

Sends a text message to a room.

- **Action**: `send-message`
- **Payload**:
  ```json
  {
    "action": "send-message",
    "message": "Hello, world!",
    "target": {
      "id": "room-id",
      "name": "Room Name"
    }
  }
  ```

#### send-audio-message

Sends an audio message to a room. The audio data should be base64 encoded.

- **Action**: `send-audio-message`
- **Payload**:
  ```json
  {
    "action": "send-audio-message",
    "message": "base64-encoded-audio-data",
    "target": {
      "id": "room-id",
      "name": "Room Name"
    }
  }
  ```

#### join-room

Joins a public room.

- **Action**: `join-room`
- **Payload**:
  ```json
  {
    "action": "join-room",
    "message": "Room Name"
  }
  ```

#### leave-room

Leaves a room.

- **Action**: `leave-room`
- **Payload**:
  ```json
  {
    "action": "leave-room",
    "message": "room-id"
  }
  ```

#### user-join

Broadcasted when a new user joins the server.

- **Action**: `user-join`

#### user-left

Broadcasted when a user leaves the server.

- **Action**: `user-left`

#### join-room-private

Joins a private room with another user.

- **Action**: `join-room-private`
- **Payload**:
  ```json
  {
    "action": "join-room-private",
    "message": "target-client-id"
  }
  ```

#### room-joined

Sent to a client when they have successfully joined a room.

- **Action**: `room-joined`
- **Payload**:
  ```json
  {
    "action": "room-joined",
    "target": {
      "id": "room-id",
      "name": "Room Name"
    }
  }
  ```

#### typing-action

Indicates that a user is typing.

- **Action**: `typing-action`
- **Payload**:
  ```json
  {
    "action": "typing-action",
    "message": "true"
  }
  ```

#### user-logged-in

Sent to a client when they have successfully logged in.

- **Action**: `user-logged-in`

#### delete-room

Deletes a room.

- **Action**: `delete-room`
- **Payload**:
  ```json
  {
    "action": "delete-room",
    "target": {
      "id": "room-id"
    }
  }
  ```

## Project Structure

```
.
├── .github/
│   └── workflows/
│       └── prod.yaml
├── .gitignore
├── .vscode/
│   ├── launch.json
│   └── tasks.json
├── chatServer.go
├── chatServer_test.go
├── client.go
├── client_test.go
├── go.mod
├── go.sum
├── main.go
├── message.go
├── message_test.go
├── room.go
└── room_test.go
```

- **`main.go`**: The entry point of the application.
- **`chatServer.go`**: Manages WebSocket connections, clients, and rooms.
- **`client.go`**: Represents a WebSocket client.
- **`room.go`**: Represents a chat room.
- **`message.go`**: Defines the message structures for WebSocket communication.
- **`*_test.go`**: Contains tests for the corresponding source files.

## Contributing

Contributions are welcome! Please feel free to submit a pull request.

## License

This project is licensed under the MIT License.
