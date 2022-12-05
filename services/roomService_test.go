package services_test

import (
	"reflect"
	"testing"

	"github.com/anjolaoluwaakindipe/videoapp/entities"
	"github.com/anjolaoluwaakindipe/videoapp/services"
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

func TestRoomService_MakeRoom(t *testing.T) {

	roomService := services.NewRoomService(&EmptyTestLogger{})

	roomService.Init();

	expected:= make(map[string][]entities.RoomParticipant)
	result := roomService.Map

	if reflect.DeepEqual(expected, result){
		t.Logf(`RoomService.Init() Passed`)
	}else{
		t.Errorf(`RoomService.Init Failed, expected "%v" but got "%v"`, expected, result )
	}

}
