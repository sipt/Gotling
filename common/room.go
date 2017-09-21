package common

import "github.com/sipt/Gotling/utils"

//IRoom 房间内部管理
type IRoom interface {
	//进入房间
	Enter(*User) error
	//高开房间
	Away(*User) error
	//把[]*User赶出房间
	DriveAway(*User, []*User) error
	//清理房间
	Clear() error
	//校验本赶走不合规矩的人，如：客户端断开连接
	CheckAndDriveAway() error
}

//Room 房间
type Room struct {
	ID    int64     //房间id
	Name  string    //房间名称
	Users *UserList //房间用户
	User  *User     //房主
}

//NewRoom 创建一个房间
func NewRoom(user *User) *Room {
	users := NewUserList()
	// users.Add(user)
	return &Room{
		ID:    utils.UniqueID(),
		User:  user,
		Users: users,
	}
}

//Enter 进入房间
func (r *Room) Enter(user *User) error {
	r.Users.Add(user)
	return nil
}

//Away 高开房间
func (r *Room) Away(user *User) error {
	r.Users.Remove(user)
	return nil
}

//DriveAway 把[]*User赶出房间
func (r *Room) DriveAway(*User, []*User) error {
	return nil
}

//Clear 清理房间
func (r *Room) Clear() error {
	return nil
}

//CheckAndDriveAway 校验本赶走不合规矩的人，如：客户端断开连接
func (r *Room) CheckAndDriveAway() error {
	return nil
}
