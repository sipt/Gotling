package transport

import (
	"net"

	"github.com/sipt/Gotling/common"
	"github.com/sipt/Gotling/encoding"
	"github.com/sipt/Gotling/network"
)

const (
	ChanBufferSize = 1 << 9
)

//ITransporter 传输
type ITransporter interface {
	Send(*common.Message) error
	Distribute(*common.Message, *common.UserNode) error
	Listen() (<-chan *MessageEntity, chan<- *MessageEntity, error)
	Receive() (*common.Message, error)
}

//MessageEntity message包
type MessageEntity struct {
	Message *common.Message
	Addr    net.Addr
	Err     error
	ErrChan chan error
}

//ServerConfig 服务器配置
type ServerConfig struct {
	IP       string
	Port     int
	Protocal string
	Addr     net.Addr
}

//NewServerConfig 初始化
func NewServerConfig(protocal, ip string, port int) *ServerConfig {
	config := &ServerConfig{
		IP:       ip,
		Port:     port,
		Protocal: protocal,
	}
	if protocal == "udp" {
		config.Addr = &net.UDPAddr{
			IP:   net.ParseIP(ip),
			Port: port,
		}
	} else if protocal == "tcp" {
		config.Addr = &net.TCPAddr{
			IP:   net.ParseIP(ip),
			Port: port,
		}
	}
	return config
}

//NewTransporter 初始化
func NewTransporter(conn network.IConn, encoder encoding.IEncoder, config *ServerConfig) ITransporter {
	return &Transporter{
		conn:         conn,
		encoder:      encoder,
		serverConfig: config,
	}
}

//Transporter 传输
type Transporter struct {
	conn         network.IConn
	encoder      encoding.IEncoder
	serverConfig *ServerConfig
	sender       chan *MessageEntity
	receiver     chan *MessageEntity
}

//Send 发送数据到服务端
func (t *Transporter) Send(msg *common.Message) error {
	entity := &MessageEntity{
		ErrChan: make(chan error, 1),
		Message: msg,
		Addr:    t.serverConfig.Addr,
	}
	t.sender <- entity
	return <-entity.ErrChan
}

//Distribute 发送数据到客户端
func (t *Transporter) Distribute(msg *common.Message, node *common.UserNode) error {
	for ; node != nil; node = node.Next() {
		user := node.Value()
		entity := &MessageEntity{
			Message: msg,
			Addr:    user.Addr,
			ErrChan: make(chan error, 1),
		}
		t.sender <- entity
		if err := <-entity.ErrChan; err != nil {
			return err
		}
	}
	return nil
}

//Receive 接收
func (t *Transporter) Receive() (*common.Message, error) {
	entity := <-t.receiver
	return entity.Message, entity.Err
}

//Listen 监听端口
func (t *Transporter) Listen() (<-chan *MessageEntity, chan<- *MessageEntity, error) {
	receiver, sender, err := t.conn.InitConn()
	if err != nil {
		return nil, nil, err
	}
	t.receiver = make(chan *MessageEntity, ChanBufferSize)
	t.sender = make(chan *MessageEntity, ChanBufferSize)
	go func() {
		for {
			select {
			case entity := <-receiver:
				if entity.Err != nil {
					t.receiver <- &MessageEntity{
						Err: entity.Err,
					}
				}
				msg, err := t.encoder.Decode(entity.Data)
				t.receiver <- &MessageEntity{
					Message: msg,
					Err:     err,
				}
			}
		}
	}()
	go func() {
		for {
			select {
			case entity := <-t.sender:
				bytes, err := t.encoder.Encode(entity.Message)
				if err != nil {
					entity.ErrChan <- err
				}
				netEntity := &network.NetEntity{
					Data:    bytes,
					Addr:    entity.Addr,
					ErrChan: make(chan error, 1),
				}
				sender <- netEntity
				entity.ErrChan <- <-netEntity.ErrChan
			}
		}
	}()
	return nil, nil, nil
}
