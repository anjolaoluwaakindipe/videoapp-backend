package services

import (
	"math/rand"
	"sync"
	"time"

	"github.com/anjolaoluwaakindipe/videoapp/entities"
	"github.com/anjolaoluwaakindipe/videoapp/utils/logger"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type RoomService struct {
	logger logger.Logger
	Mutex  sync.RWMutex
	Map    map[string][]entities.RoomParticipant
}

// Init: initializes the room structure
func (rs *RoomService) Init() {
	rs.Map = make(map[string][]entities.RoomParticipant)
}

// Get: returns the array of participants in the room
func (rs *RoomService) Get(roomID string) []entities.RoomParticipant {
	rs.Mutex.Lock()
	defer rs.Mutex.Unlock()

	return rs.Map[roomID]
}

// CreateRoom: generates a unique ID and insert in the hashmap then return the id
func (rs *RoomService) CreateRoom() string {
	// generate a unique ID and insert in the hashmap
	rs.Mutex.Lock()
	defer rs.Mutex.Unlock()

	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 8)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	roomId := string(b)
	rs.Map[roomId] = []entities.RoomParticipant{}

	return roomId
}

// InsertIntoRoom: creates a participant and add it into room structure
func (rs *RoomService) InsertIntoRoom(roomID string, host bool, conn *websocket.Conn) string {
	rs.Mutex.Lock()
	defer rs.Mutex.Unlock()

	newParicticipantId := uuid.New()

	p := entities.RoomParticipant{Host: host, Conn: conn, Id: newParicticipantId.String()}

	rs.logger.Info("Inserting into Room with room Id: " + roomID)
	rs.Map[roomID] = append(rs.Map[roomID], p)

	rs.logger.Infof("Participants in %v: %v", roomID ,len(rs.Map[roomID]))
	return newParicticipantId.String()
}

// CloseConnection closes a connection and removes the connection from a room
func (rs *RoomService) CloseConnection(conn *websocket.Conn, roomID string) {
	participants := rs.Map[roomID]

	for i := 0; i < len(participants); i++ {
		if participants[i].Conn == conn {
			participants[i] = participants[len(participants)-1]
			rs.Map[roomID] = participants[:len(participants)-1]
			rs.logger.Infof("removed participant from room %v : %v", roomID, rs.Map[roomID])
			break
		} else {
			continue
		}
	}

	conn.Close()
}

func (rs *RoomService) DeleteRoom(roomID string) {
	rs.Mutex.Lock()
	defer rs.Mutex.Unlock()

	delete(rs.Map, roomID)
}

func (rs *RoomService) ShowMap() map[string][]entities.RoomParticipant {
	return rs.Map
}

// RoomService Constructor
func NewRoomService(logger logger.Logger) *RoomService {
	var roomService RoomService = RoomService{logger: logger}
	roomService.Init()
	return &roomService
}
