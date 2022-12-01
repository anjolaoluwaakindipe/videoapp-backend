package services

import (
	"github.com/anjolaoluwaakindipe/videoapp/entities"
	"github.com/gorilla/websocket"
)

type RoomServiceInterface interface {
	Init()
	Get(roomID string) []entities.RoomParticipant
	CreateRoom() string
	InsertIntoRoom(roomID string, host bool, coon *websocket.Conn)
	DeleteRoom(roomID string)
	ShowMap() map[string][]entities.RoomParticipant
	CloseConnection(conn *websocket.Conn, roomID string)
}
