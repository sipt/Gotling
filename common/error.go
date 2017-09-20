package common

import "fmt"

//Error 自定义error
type Error struct {
	Code string
	Msg  string
}

//Error error to string
func (e *Error) Error() string {
	return fmt.Sprintf("%s(%s)", e.Msg, e.Code)
}
