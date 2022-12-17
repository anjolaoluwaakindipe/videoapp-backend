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

			// log out what room was created
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
		// wait for broadcast message
		msg := <-*broadcast
		
		// broadcast to all other participant except self
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
			// get room id from url params
			roomID := r.URL.Query().Get("roomID")

			// trim whitespace from roomID id: not sure I need to do this 
			if strings.TrimSpace(roomID) == "" {
				rh.logger.Infoln("roomID missing in URL Parameeters")
				err := errs.NewBadRequestError("RoomID query parameter invalid")
				res(rw, err, err.Code)
				return
			}

			// check if room Id is part of map
			if _, ok := rh.roomService.ShowMap()[roomID]; !ok  {
				err := errs.NewBadRequestError("Room ID is invalid")
				res(rw, err, err.Code)
				return
			}

			// create a websocket upgrader
			var upgrader = websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			}

			// upgrade http connection to websocket connection 
			ws, socketErr := upgrader.Upgrade(rw, r, nil)
			
			// check if there was an error in upgrading connection
			if socketErr != nil {
				rh.logger.Fatal("Web socket upgrade err: " + socketErr.Error())
				err := errs.NewInternalServerError()
				res(rw, err, err.Code)
				return 
			}

			// defer the closeConnection method if websocket conection is closed intentionally or not
			defer rh.roomService.CloseConnection(ws, roomID)

			// TODO: make a participant a host

			// insert websocket into connection room (via RoomService)
			rh.roomService.InsertIntoRoom(roomID, false, ws)

			// make a broadcast channel
			broadcast := make(chan entities.BroadcastMsg)

			// run broadcaster method concurrently before infinite for loop
			go rh._broadcaster(&broadcast)

			// run infinite for loop to detect messages 
			for {
				// variable to hold  broadcast Msg
				var msg entities.BroadcastMsg

				// read the next message sent by particpant
				err := ws.ReadJSON(&msg.Message)
				
				// if there was an error disconnect user from room
				if err != nil {
					rh.logger.Info("Web socket json read err: " + err.Error())
					break;
				}

				// set mesage client and room Id
				msg.Client = ws
				msg.RoomID = roomID

				// broadcast participant's message to other participants through broad cast channel
				broadcast <- msg
			}
		},
	}
}

// RoomHandler Constructor
func NewRoomHandler(roomService services.RoomServiceInterface, logger logger.Logger) RoomHandler {
	return RoomHandler{roomService: roomService, logger: logger}
}
