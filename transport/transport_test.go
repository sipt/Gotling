package transport

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/sipt/Gotling/common"
	"github.com/sipt/Gotling/encoding"
	"github.com/sipt/Gotling/network"
	"github.com/sipt/Gotling/room"
)

func TestTransport(t *testing.T) {
	data := "这是一条测试数据"
	testMsg := &common.Message{
		UserID:    int64(1),
		RoomID:    int64(2),
		Type:      int16(3),     //消息类型
		Data:      []byte(data), //消息内容
		Timescanp: time.Now(),
	}

	encoder := encoding.NewJSONEncoder()
	conn := network.NewUDPConn(20010)
	config := NewServerConfig("udp", "127.0.0.1", 20010)
	transporter := NewTransporter(conn, encoder, config)
	building := room.NewBuilding()
	_, _, err := transporter.Listen()
	if err != nil {
		t.Errorf("listen failed, err:%v", err)
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		var msg *common.Message
		var err error
		for {
			msg, _, err = transporter.Receive()
			if err != nil {
				t.Errorf("receive msg error, err:%v", err)
			}
			if msg.UserID != testMsg.UserID || msg.RoomID != testMsg.RoomID || msg.Timescanp.Nanosecond() != testMsg.Timescanp.Nanosecond() || string(msg.Data) != string(testMsg.Data) || msg.Type != testMsg.Type {
				t.Fail()
			}
			wg.Done()
		}
	}()
	err = transporter.Send(testMsg)
	if err != nil {
		t.Errorf("Send message failed, err:%v", err)
	}
	user := &common.User{
		Addr: &net.UDPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 20010,
		},
	}
	room, _ := building.CheckIn(user)
	err = transporter.Distribute(testMsg, room.Users.Head())
	if err != nil {
		t.Errorf("Distribute message failed, err:%v", err)
	}
	wg.Wait()
}
