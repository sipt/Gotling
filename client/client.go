package main

import (
	"flag"
	"fmt"

	"github.com/sipt/Gotling/common"
	"github.com/sipt/Gotling/encoding"
	"github.com/sipt/Gotling/exec"
	"github.com/sipt/Gotling/network"
	"github.com/sipt/Gotling/room"
	"github.com/sipt/Gotling/transport"
	"github.com/sipt/Gotling/utils"
)

func main() {
	var cport = flag.Int("cport", 20011, "udp listen port")
	var sport = flag.Int("sport", 20010, "udp listen port")
	var ip = flag.String("ip", "127.0.0.1", "udp server ip")
	flag.Parse()
	encoder := encoding.NewJSONEncoder()
	conn := network.NewUDPConn(*cport)
	config := transport.NewServerConfig("udp", *ip, *sport)
	transporter := transport.NewTransporter(conn, encoder, config)
	building := room.NewBuilding()
	executor := exec.NewExecutor(exec.ClientMode, transporter, building)
	result := executor.LoopExec()
	var commond string
	var id int64
	var user = &common.User{
		ID: utils.UniqueID(),
	}
	var roomID int64 = -1
	go func() {
		for {
			msg := <-result
			fmt.Printf("\n[%d]:%s", msg.UserID, string(msg.Data))
			if roomID < 0 {
				fmt.Print("\n➜ ")
			} else {
				fmt.Printf("\n[%d]➜ ", roomID)
			}
		}
	}()
	fmt.Printf("Gotaling [%d]:\n请输入昵称：\n➜ ", user.ID)
	fmt.Scanf("%s", &user.Nick)
	for {
		if roomID < 0 {
			fmt.Print("\n➜ ")
		} else {
			fmt.Printf("\n[%d]➜ ", roomID)
		}
		fmt.Scanf("%s %d", &commond, &id)
		fmt.Println(commond)
		if commond == "" {
			if roomID >= 0 {
				msg := &common.Message{
					RoomID: roomID,
					UserID: user.ID,
					Type:   common.MsgTypeNil,
				}
				err := executor.ToServer(msg)
				if err != nil {
					fmt.Println("err:", err)
				}
			}
			continue
		}
		switch commond {
		case "into":
			msg := &common.Message{
				RoomID: id,
				UserID: user.ID,
				Type:   common.MsgTypeUserIntoRoom,
			}
			err := executor.ToServer(msg)
			if err != nil {
				fmt.Println("err:", err)
			} else {
				roomID = id
			}
		case "out":
			msg := &common.Message{
				RoomID: roomID,
				UserID: user.ID,
				Type:   common.MsgTypeUserAwayRoom,
			}
			err := executor.ToServer(msg)
			if err != nil {
				fmt.Println("err:", err)
			} else {
				roomID = -1
			}
		default:
			msg := &common.Message{
				Data:   []byte(commond),
				RoomID: roomID,
				UserID: user.ID,
				Type:   common.MsgType,
			}
			err := executor.ToServer(msg)
			if err != nil {
				fmt.Println("err:", err)
			}
		}
		commond = ""
	}
}
