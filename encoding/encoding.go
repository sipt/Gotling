package encoding

import (
	"encoding/json"

	"github.com/sipt/Gotling/common"
)

//IEncoder 编码解码器
type IEncoder interface {
	//Encode 编码
	Encode(*common.Message) ([]byte, error)
	//Decode 解码
	Decode([]byte) (*common.Message, error)
}

//NewJSONEncoder 初始化
func NewJSONEncoder() IEncoder {
	return new(JSONEncoder)
}

//JSONEncoder json encoder
type JSONEncoder struct {
}

//Encode 编码
func (e *JSONEncoder) Encode(msg *common.Message) ([]byte, error) {
	return json.Marshal(msg)
}

//Decode 解码
func (e *JSONEncoder) Decode(bytes []byte) (*common.Message, error) {
	msg := &common.Message{}
	err := json.Unmarshal(bytes, msg)
	return msg, err
}
