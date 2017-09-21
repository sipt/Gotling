package room

import (
	"sync"

	"github.com/sipt/Gotling/common"
)

const (
	//DefaultSize 默认长度
	DefaultSize = 1 << 6
)

//IBuilding 负责房间注册和注销
type IBuilding interface {
	//注册房间
	CheckIn(*common.User) (*common.Room, error)
	//回收房间
	CheckOut(*common.Room, *common.User) error
	Rooms() []*common.Room
	Room(id int64) (*common.Room, error)
}

//NewBuilding 初始化
func NewBuilding() IBuilding {
	return &Building{
		rooms: make([]*common.Room, 0, DefaultSize),
	}
}

//Building 楼
type Building struct {
	sync.RWMutex
	rooms []*common.Room
}

//CheckIn 注册房间
func (b *Building) CheckIn(user *common.User) (*common.Room, error) {
	b.Lock()
	defer b.Unlock()
	room := common.NewRoom(user)
	b.rooms = append(b.rooms, room)
	return room, nil
}

//CheckOut 回收房间
func (b *Building) CheckOut(room *common.Room, user *common.User) error {
	b.Lock()
	defer b.Unlock()
	for i, r := range b.rooms {
		if r.ID == room.ID {
			if room.User == nil || user.ID == room.User.ID {
				b.rooms = append(b.rooms[:i], b.rooms[i+1:]...)
			} else {
				return &common.Error{Code: "BUILD001", Msg: "权限不足"}
			}
		}
	}
	return nil
}

//Rooms Rooms列表
func (b *Building) Rooms() []*common.Room {
	b.RLock()
	defer b.RUnlock()
	return b.rooms
}

//Room 通过Room.ID获取Room
func (b *Building) Room(id int64) (*common.Room, error) {
	b.RLock()
	defer b.RUnlock()
	for _, r := range b.rooms {
		if r.ID == id {
			return r, nil
		}
	}
	return nil, &common.Error{
		Code: "BUILD002",
		Msg:  "找不到对应房间",
	}
}
