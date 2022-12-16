package entities

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type RoomParticipant struct {
	Host bool
	Conn *websocket.Conn
}

func (rp RoomParticipant) String() string {
	return fmt.Sprintf(`{Host: %v, Conn: %v}`, rp.Host, rp.Conn)
} 