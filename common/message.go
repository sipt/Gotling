package common

import "time"

const (
	//无操作
	MsgTypeNil          = 0
	MsgType             = 1
	MsgTypeRoomAdd      = 100
	MsgTypeRoomDel      = 110
	MsgTypeRoomList     = 120
	MsgTypeRoomUpdate   = 130
	MsgTypeUserAdd      = 200
	MsgTypeUserAwayRoom = 210
	MsgTypeUserIntoRoom = 211
)

//Message 信息体
type Message struct {
	Type      int16     //消息类型
	Data      []byte    //消息内容
	RoomID    int64     //房间ID
	UserID    int64     //发送者ID
	Timescanp time.Time //时间戳
}
