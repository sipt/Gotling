package room

import (
	"github.com/sipt/Gotling/common"
	"github.com/sipt/Gotling/encoding"
)

//ISpeaker 发言
type ISpeaker interface {
	Speak(*common.User, *common.Room) error
}

//Speaker 发言
type Speaker struct {
	encoder encoding.IEncoder
}
