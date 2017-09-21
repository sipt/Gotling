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
)

const ()

func main() {
	var port = flag.Int("port", 20010, "udp listen port")
	encoder := encoding.NewJSONEncoder()
	conn := network.NewUDPConn(*port)
	transporter := transport.NewTransporter(conn, encoder, nil)
	building := room.NewBuilding()
	executor := exec.NewExecutor(exec.ServerMode, transporter, building)
	executor.LoopExec()
	fmt.Println("Gotaling:")
	var commond string
	var id int64
	for {
		fmt.Print("\n➜ ")
		fmt.Scanf("%s %d", &commond, &id)
		switch commond {
		case "checkin":
			room, err := building.CheckIn(nil)
			if err != nil {
				fmt.Println("err:", err)
			} else {
				fmt.Println("roomID:", room.ID)
			}
		case "checkout":
			err := building.CheckOut(&common.Room{ID: id}, nil)
			if err != nil {
				fmt.Println("err:", err)
			} else {
				fmt.Println("check out success")
			}
		case "rooms":
			rooms := building.Rooms()
			for _, room := range rooms {
				fmt.Println(room.ID)
			}
		case "room":
			room, err := building.Room(id)
			if err != nil {
				fmt.Println("err:", err)
			} else {
				node := room.Users.Head()
				fmt.Println("user list: ")
				for ; node != nil; node = node.Next() {
					fmt.Println(node.Value().ID)
				}
			}
		}
	}
}
