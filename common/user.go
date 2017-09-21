package common

import (
	"net"
	"sync"
)

const (
	//TCPProtocal tcp 协议常量
	TCPProtocal = "tcp"
	//UDPProtocal udp 协议常量
	UDPProtocal = "udp"
)

//User 用户
type User struct {
	ID       int64  //用户ID
	Account  string //用户帐号
	Password string //用户密码
	Nick     string //用户昵称

	Addr net.Addr
}

//IAddr 地址
type IAddr interface {
	//return ip:port
	Addr() string
	GetIP() string
	GetPort() string
	//协议 ["tcp","udp"]
	Protocal() string
}

//NetInfo 网络来源信息
type NetInfo struct {
	ip       string //IP地址
	port     string //端口号
	protocal string //协议
}

//Addr return ip:port
func (n *NetInfo) Addr() string {
	return n.ip + ":" + n.port
}

//GetIP return ip
func (n *NetInfo) GetIP() string {
	return n.ip
}

//GetPort return port
func (n *NetInfo) GetPort() string {
	return n.port
}

//Protocal 协议 ["tcp","udp"]
func (n *NetInfo) Protocal() string {
	return n.protocal
}

//UserNode 用户节点
type UserNode struct {
	user       *User     //当前节点数据
	prev, next *UserNode //上个节点、下个节点数据
}

//Prev 上个节点
func (n *UserNode) Prev() *UserNode {
	return n.prev
}

//Next 下个节点
func (n *UserNode) Next() *UserNode {
	return n.next
}

//Append 添加节点
func (n *UserNode) Append(user *User) *UserNode {
	node := &UserNode{
		user: user,
		next: nil,
		prev: n,
	}
	n.next = node
	return node
}

//Remove 删除本节点
func (n *UserNode) Remove() {
	n.prev = n.next
}

//Value 返回用户
func (n *UserNode) Value() *User {
	return n.user
}

//UserList 用户链表
type UserList struct {
	sync.RWMutex
	head, tail *UserNode //链表头、尾、当前节点
	len        int64     //链表长
}

//Head 返回头节点
func (l *UserList) Head() *UserNode {
	l.RLock()
	defer l.RUnlock()
	return l.head
}

//Tail 返回尾节点
func (l *UserList) Tail() *UserNode {
	l.RLock()
	defer l.RUnlock()
	return l.tail
}

//Add 添加节点
func (l *UserList) Add(user *User) *UserNode {
	l.Lock()
	defer l.Unlock()
	var node *UserNode
	if l.len == 0 {
		node = &UserNode{
			user: user,
			next: nil,
			prev: nil,
		}
		l.head = node
		l.tail = node
	} else {
		node = l.tail.Append(user)
		l.tail = node
	}
	l.len++
	return node
}

//Remove 删除节点
func (l *UserList) Remove(user *User) *UserNode {
	l.Lock()
	defer l.Unlock()
	for node := l.Head(); node != nil; node = node.next {
		if node.Value().ID == user.ID {
			node.Remove()
			l.len--
			return node
		}
	}
	return nil
}
