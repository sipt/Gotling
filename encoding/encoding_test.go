package encoding

import (
	"testing"
	"time"

	"github.com/sipt/Gotling/common"
)

func TestJSONEncoding(t *testing.T) {
	now := time.Now()
	msg := &common.Message{
		UserID:    int64(1),
		RoomID:    int64(2),
		Type:      int8(3),         //消息类型
		Data:      []byte("我的一句话"), //消息内容
		Timescanp: now,
	}
	var encoder IEncoder = new(JSONEncoder)
	data, err := encoder.Encode(msg)
	if err != nil {
		t.Errorf("encode failed, err:%v", err)
	}
	msg2, err := encoder.Decode(data)
	if err != nil {
		t.Errorf("decode failed, err:%v", err)
	}
	if msg.UserID != msg2.UserID || msg.RoomID != msg2.RoomID || msg.Timescanp.Nanosecond() != msg2.Timescanp.Nanosecond() || string(msg.Data) != string(msg2.Data) || msg.Type != msg2.Type {
		t.Error("msg != msg2")
	}
}
