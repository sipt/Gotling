package common

import "testing"

func TestUserList(t *testing.T) {
	users := NewUserList()
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
		t.Errorf("value error, value=%v", node.Value().ID)
	}
	node = node.Next()
	if node == nil || node.Value().ID != 234234 {
		t.Errorf("value error, value=%v", node.Value().ID)
	}
	node = nil
	users.Remove(&User{
		ID: 123123,
	})
	if users.len != 1 {
		t.Error("value length error")
	}
	node = users.Head()
	if node == nil || node.Value().ID != 234234 {
		t.Errorf("value error, value=%v", node.Value().ID)
	}
}
