package room

import (
	"testing"

	"github.com/sipt/Gotling/common"
	"github.com/sipt/Gotling/utils"
)

func TestCheckInAndCheckOut(t *testing.T) {
	building := NewBuilding()
	user := &common.User{ID: utils.UniqueID()}
	room, err := building.CheckIn(user)
	if err != nil {
		t.Errorf("check in failed, err:%v", err)
	}
	rooms := building.Rooms()
	if len(rooms) != 1 {
		t.Errorf("check rooms length failed, err:%v", err)
	}
	temp, err := building.Room(room.ID)
	if err != nil {
		t.Errorf("room not found, err:%v", err)
	} else if temp.ID != room.ID {
		t.Errorf("found room failed, err:%v", err)
	}
	err = building.CheckOut(room, user)
	if err != nil {
		t.Errorf("check out failed, err:%v", err)
	}
	rooms = building.Rooms()
	if len(rooms) != 0 {
		t.Errorf("check rooms length failed, err:%v", err)
	}
}
