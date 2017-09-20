package common

import "time"

//Message 信息体
type Message struct {
	Type      int8      //消息类型
	Data      []byte    //消息内容
	RoomID    int64     //房间ID
	UserID    int64     //发送者ID
	Timescanp time.Time //时间戳
}
