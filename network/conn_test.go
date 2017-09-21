package network

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestUDPConn(t *testing.T) {
	data := "这是一条测试数据"
	count := 10
	ip := "127.0.0.1" //"120.26.51.19"
	port := 20010
	var conn IConn = &UDPConn{port: port}
	receiver, sender, err := conn.InitConn()
	if err != nil {
		t.Errorf("listen port %d failed, err: %v", port, err)
	}
	var wg sync.WaitGroup
	wg.Add(count)
	go func() {
		var bytes []byte
		for index := 0; index < count; index++ {
			entity := <-receiver
			if entity.Err != nil {
				t.Errorf("receive data failed, err:%v", err)
			}
			bytes = entity.Data
			if string(bytes) != data {
				fmt.Printf("'%s' != '%s' -> %v\n", string(bytes), data, string(bytes) != data)
				t.Fail()
			}
			wg.Done()
		}
	}()
	for index := 0; index < count; index++ {
		entity := &NetEntity{
			Data: []byte(data),
			Addr: &net.UDPAddr{
				IP:   net.ParseIP(ip),
				Port: port,
			},
			ErrChan: make(chan error, 1),
		}
		sender <- entity
		err := <-entity.ErrChan
		if err != nil {
			t.Errorf("send data failed, err:%v", err)
		}
	}
	wg.Wait()
}
