package exec

import (
	"github.com/sipt/Gotling/common"
	"github.com/sipt/Gotling/room"
	"github.com/sipt/Gotling/transport"
)

const (
	ServerMode int8 = 0
	ClientMode int8 = 1
)

//IExecutor 执行者
type IExecutor interface {
	ToServer(*common.Message) error
	Distribute(*common.Message) error
	LoopExec() <-chan *common.Message
}

//NewExecutor 初始化
func NewExecutor(mode int8, t transport.ITransporter, b room.IBuilding) IExecutor {
	return &Executor{
		transporter: t,
		Building:    b,
		mode:        mode,
	}
}

//Executor 发言
type Executor struct {
	transporter transport.ITransporter
	Building    room.IBuilding
	mode        int8
}

//ToServer 对（服务器）说话
func (e *Executor) ToServer(msg *common.Message) error {
	return e.transporter.Send(msg)
}

//Distribute （服务器）扩散消息
func (e *Executor) Distribute(msg *common.Message) error {
	room, err := e.Building.Room(msg.RoomID)
	if err != nil {
		return err
	}
	return e.transporter.Distribute(msg, room.Users.Head())
}

//LoopExec 循环扩散消息
func (e *Executor) LoopExec() <-chan *common.Message {
	e.transporter.Listen()
	var exec func() <-chan *common.Message
	switch e.mode {
	case ServerMode:
		exec = e.serverExec
	case ClientMode:
		exec = e.clientExec
	}
	return exec()
}

func (e *Executor) serverExec() <-chan *common.Message {
	go func() {
		for {
			msg, addr, err := e.transporter.Receive()
			if err != nil {
				//TODO 错误处理
			} else {
				switch msg.Type {
				case common.MsgTypeNil:
				case common.MsgType:
					err = e.Distribute(msg)
				case common.MsgTypeRoomAdd:
				case common.MsgTypeRoomDel:
				case common.MsgTypeRoomList:
				case common.MsgTypeRoomUpdate:
				case common.MsgTypeUserAdd:
				case common.MsgTypeUserAwayRoom:
					room, err := e.Building.Room(msg.RoomID)
					if err == nil {
						room.Users.Remove(&common.User{ID: msg.UserID})
					}
				case common.MsgTypeUserIntoRoom:
					room, err := e.Building.Room(msg.RoomID)
					if err == nil {
						room.Users.Add(&common.User{ID: msg.UserID, Addr: addr})
					}
				default:
					err = &common.Error{Code: "EXE001", Msg: "op type not supported"}
				}
				if err != nil {
					//TODO 错误处理
				}
			}
		}
	}()
	return nil
}
func (e *Executor) clientExec() <-chan *common.Message {
	result := make(chan *common.Message, 64)
	go func() {
		for {
			msg, _, err := e.transporter.Receive()
			if err != nil {
				//TODO 错误处理
			} else {
				switch msg.Type {
				case common.MsgTypeNil:
				case common.MsgType:
					result <- msg
				default:
					err = &common.Error{Code: "EXE001", Msg: "op type not supported"}
				}
				if err != nil {
					//TODO 错误处理
				}
			}
		}
	}()
	return result
}
