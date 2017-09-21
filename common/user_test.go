package common

import "testing"

func TestUserList(t *testing.T) {
	users := new(UserList)
	users.Add(&User{
		ID: 123123,
	})
	users.Add(&User{
		ID: 234234,
	})
	if users.len != 2 {
		t.Error("value length error")
	}
	node := users.Head()
	if node == nil || node.Value().ID != 123123 {
		t.Error("value error")
	}
	node = node.Next()
	if node == nil || node.Value().ID != 234234 {
		t.Error("value error")
	}
}
