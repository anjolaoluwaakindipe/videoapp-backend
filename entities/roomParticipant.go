package entities

import "github.com/gorilla/websocket"

type RoomParticipant struct {
	Host bool
	Conn *websocket.Conn
}