package network

import (
	"net"

	"github.com/sipt/Gotling/common"
)

const (
	//ChannelBufferSize 通道缓冲区大小
	ChannelBufferSize = 1 << 9
	//DatagramSize 数据报大小
	DatagramSize = 1 << 9
)

//NetEntity 网络包
type NetEntity struct {
	Addr    net.Addr
	Data    []byte
	ErrChan chan error
}

//NewNetEntity 初始化一个netEntity
func NewNetEntity(addr net.Addr, data []byte, err error) *NetEntity {
	if err != nil {
		return &NetEntity{
			ErrChan: make(chan error, 1),
		}
	}
	return &NetEntity{
		Addr: addr,
		Data: data,
	}
}

//IConn 连接
type IConn interface {
	//初始化网络
	InitConn(port int) (<-chan *NetEntity, chan<- *NetEntity, error)
	Close()
}

//UDPConn udp连接
type UDPConn struct {
	conn *net.UDPConn
}

//InitConn 监听端口
func (c *UDPConn) InitConn(port int) (<-chan *NetEntity, chan<- *NetEntity, error) {
	conn, err := net.ListenUDP(common.UDPProtocal, &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: port,
	})
	if err != nil {
		return nil, nil, err
	}
	receiver := make(chan *NetEntity, ChannelBufferSize)
	sender := make(chan *NetEntity, ChannelBufferSize)
	go func() {
		for {
			data := make([]byte, DatagramSize)
			length, udpAddr, err := conn.ReadFromUDP(data)
			receiver <- NewNetEntity(udpAddr, data[:length], err)
		}
	}()
	go func() {
		for {
			entity := <-sender
			udpAddr, ok := entity.Addr.(*net.UDPAddr)
			if !ok {
				entity.ErrChan <- &common.Error{Code: "NET001", Msg: "地址信息非UDPAddr"}
			}
			length, err := conn.WriteToUDP(entity.Data, udpAddr)
			if err != nil {
				entity.ErrChan <- err
			} else if length <= 0 {
				entity.ErrChan <- &common.Error{Code: "NET002", Msg: "发送数据为空"}
			} else {
				entity.ErrChan <- nil
			}
		}
	}()
	return receiver, sender, nil
}

//Close close conn
func (u *UDPConn) Close() {
	u.conn.Close()
}
