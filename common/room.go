package common

//Room 房间
type Room struct {
	ID    int64     //房间id
	Name  string    //房间名称
	Users *UserList //房间用户
	User  *User     //房主
}
