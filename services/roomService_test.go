package services_test

import (
	"reflect"
	"testing"

	"github.com/anjolaoluwaakindipe/videoapp/entities"
	"github.com/anjolaoluwaakindipe/videoapp/services"
	"github.com/gorilla/websocket"
)

type EmptyTestLogger struct{}

func (zl *EmptyTestLogger) Info(message string) {
}

func (zl *EmptyTestLogger) Infoln(message string) {
}

func (zl *EmptyTestLogger) Infof(template string, args ...interface{}) {
}

func (zl *EmptyTestLogger) Debug(message string) {
}

func (zl *EmptyTestLogger) Error(message string) {
}

func (zl *EmptyTestLogger) Fatal(message string) {
}

func (z1 *EmptyTestLogger) Fatalf(template string, args ...interface{}) {
}

func (zl *EmptyTestLogger) Fatalln(message string) {
}

// new participant struct for holding inputs
type NewParticipant struct {
	roomId string
	host   bool
	conn   *websocket.Conn
}

func InitRoomService() *services.RoomService {
	roomService := services.NewRoomService(&EmptyTestLogger{})

	roomService.Init()
	return roomService
}

// INITIATE THE ROOMSERVICE TEST
func TestRoomService_Init(t *testing.T) {

	roomService := InitRoomService()
	expected := make(map[string][]entities.RoomParticipant)
	result := roomService.Map

	if reflect.DeepEqual(expected, result) {
		t.Logf(`RoomService.Init() PASSED`)
	} else {
		t.Errorf(`RoomService.Init FAILED, expected "%v" but got "%v"`, expected, result)
	}

}

// CREATE ROOM TEST
func TestRoomService_CreateRoom(t *testing.T) {
	roomService := InitRoomService()

	var roomId string = roomService.CreateRoom()

	_, isPresent := roomService.Map[roomId]

	if isPresent == true {
		t.Logf("RoomService.CreateRoom PASSED")
	} else {
		t.Errorf(`RoomService.CreateRoom FAILED, expected "true" but got "%v"`, isPresent)
	}
}

// INSERT PARTICIPANT INTO ROOM TEST
func TestRoomService_InsertIntoRoom(t *testing.T) {

	// initiate room service
	roomService := InitRoomService()

	// create a room and get room Id
	roomId := roomService.CreateRoom()

	// create test tablle of inputs that insert a particpant
	testTable := []struct {
		name   string
		input  NewParticipant
		output int
	}{
		{name: "First participant insert", input: NewParticipant{roomId: roomId, host: false, conn: &websocket.Conn{}}, output: 1},
		{name: "Second particpant insert", input: NewParticipant{roomId: roomId, host: false, conn: &websocket.Conn{}}, output: 2},
		{name: "Third participant insert", input: NewParticipant{roomId: roomId, host: false, conn: &websocket.Conn{}}, output: 3},
	}
	// check the length of participants in the room
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			// each input
			input := test.input

			// insert each particpant into the room
			roomService.InsertIntoRoom(input.roomId, input.host, input.conn)

			// get the number of participants
			numberOfParticipants := len(roomService.Map[roomId])

			if numberOfParticipants == test.output {
				t.Logf(`RoomService.InsertIntoRoom PASSED (%v)`, test.name)
			} else {
				t.Logf(`RoomService.InsertIntoRoom FAILED, expected "%v" but got "%v"`, test.output, numberOfParticipants)
			}
		})
	}
}

// TESTS IF A ROOM WAS DELETED FROM THE ROOM SERVICE
func TestRoomService_DeleteRoom(t *testing.T) {
	// initiate
	roomService := InitRoomService()

	// create room
	var roomId string = roomService.CreateRoom()

	// delete room
	roomService.DeleteRoom(roomId)

	// check if room is not in roomService.Map
	_, isPresent := roomService.Map[roomId]

	if isPresent == false {
		t.Logf(`RoomService.DeleteRoom PASSED`)
	} else {
		t.Errorf(`RoomService.DeleteRoom FAILED, expected "%v" but got "%v"`, false, isPresent)
	}
}

// GET PARTICIPANTS FROM A SPECIFIC ROOM TEST
func TestRoomService_Get(t *testing.T) {

	// test table for get test
	testTable := []struct {
		name           string
		inputs         []NewParticipant
		expectedOutput []entities.RoomParticipant
	}{
		{name: "Test a single participant", inputs: []NewParticipant{{host: false, conn: &websocket.Conn{}}}},
		{name: "Test zero participants", inputs: make([]NewParticipant, 0), expectedOutput: make([]entities.RoomParticipant, 0)},
	}

	// run individual tests suite
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {

			// initiate room service
			roomSeervice := InitRoomService()

			// get room id
			var roomId string = roomSeervice.CreateRoom()

			//  setup room id for each input and map input to expected output
			for i := 0; i < len(test.inputs); i++ {
				// test input
				input := test.inputs[i]
				input.roomId = roomId
				roomSeervice.InsertIntoRoom(input.roomId, input.host, input.conn)
				test.expectedOutput = append(test.expectedOutput, entities.RoomParticipant{Host: input.host, Conn: input.conn})

			}

			// get result
			result := roomSeervice.Get(roomId)

			// assertion for result and expected
			if reflect.DeepEqual(result, test.expectedOutput) {
				t.Logf("RoomService.Get PASSED (%v)", test.name)
			} else {
				t.Logf(`RoomService.Get FAILED, expected "%v" but got "%v"`, test.expectedOutput, result)
			}
		})

	}

}
