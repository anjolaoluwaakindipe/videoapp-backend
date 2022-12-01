package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/anjolaoluwaakindipe/videoapp/dto/errs"
	"github.com/anjolaoluwaakindipe/videoapp/dto/response"
	"github.com/anjolaoluwaakindipe/videoapp/entities"
	"github.com/anjolaoluwaakindipe/videoapp/services"
	"github.com/anjolaoluwaakindipe/videoapp/utils/logger"
	"github.com/gorilla/websocket"
)

type RoomHandler struct {
	roomService services.RoomServiceInterface
	logger      logger.Logger
	broadcast   chan entities.BroadcastMsg
}

// handler: allows a user to create a room
func (rh RoomHandler) createRoom() RouteHandler {
	return RouteHandler{
		Path:    "/room/create",
		Methods: []string{http.MethodGet, http.MethodOptions},
		HandlerFunc: func(rw http.ResponseWriter, r *http.Request) {
			// create a room and return room id
			roomID := rh.roomService.CreateRoom()

			rh.logger.Infof("%v", rh.roomService.ShowMap())

			rw.Header().Add("Content-Type", "application/json")
			rw.Header().Set("Access-Control-Allow-Origin", "*")

			// send response
			json.NewEncoder(rw).Encode(response.CreateRoomRes{RoomID: roomID})
		},
	}
}

func (rh RoomHandler) _broadcaster(broadcast *chan entities.BroadcastMsg) {
	for {
		msg := <-*broadcast

		for _, client := range rh.roomService.ShowMap()[msg.RoomID] {
			if client.Conn != msg.Client {
				err := client.Conn.WriteJSON(msg.Message)

				if err != nil {
					rh.logger.Fatal("Error broadcasting websocket message: " + err.Error())
					client.Conn.Close()
				}
			}
		}
	}
}

// handler: allows a user to join a room
func (rh RoomHandler) joinRoom() RouteHandler {
	return RouteHandler{
		Path:    "/room/join",
		Methods: []string{"GET"},
		HandlerFunc: func(rw http.ResponseWriter, r *http.Request) {
			roomID := r.URL.Query().Get("roomID")

			if strings.TrimSpace(roomID) == "" {
				rh.logger.Infoln("roomID missing in URL Parameeters")
				err := errs.NewBadRequestError("RoomID query parameter invalid")
				res(rw, err, err.Code)
				return
			}

			if _, ok := rh.roomService.ShowMap()[roomID]; ok == false {
				err := errs.NewBadRequestError("Room ID is invalid")
				res(rw, err, err.Code)
				return
			}

			var upgrader = websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			}

			ws, socketErr := upgrader.Upgrade(rw, r, nil)
			

			if socketErr != nil {
				rh.logger.Fatal("Web socket upgrade err: " + socketErr.Error())
			}
			defer rh.roomService.CloseConnection(ws, roomID)

			rh.roomService.InsertIntoRoom(roomID, false, ws)

			broadcast := make(chan entities.BroadcastMsg)

			go rh._broadcaster(&broadcast)

			for {
				var msg entities.BroadcastMsg

				err := ws.ReadJSON(&msg.Message)

				if err != nil {
					rh.logger.Info("Web socket json read err: " + err.Error())
					break;
				}
				msg.Client = ws
				msg.RoomID = roomID
				broadcast <- msg
			}
		},
	}
}

// RoomHandler Constructor
func NewRoomHandler(roomService services.RoomServiceInterface, logger logger.Logger) RoomHandler {
	return RoomHandler{roomService: roomService, logger: logger}
}
