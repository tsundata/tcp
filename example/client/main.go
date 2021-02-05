package main

import (
	"github.com/tsundata/tcp/network"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:5678")
	if err != nil {
		panic(err)
	}

	for {
		// write
		pack := network.NewPack()
		msg, _ := pack.Pack(network.NewMessage(0, []byte(time.Now().String())))
		_, err := conn.Write(msg)
		if err != nil {
			log.Println(err)
			return
		}

		// read
		headData := make([]byte, pack.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			log.Println(err)
			break
		}
		head, err := pack.Unpack(headData)
		if err != nil {
			log.Println(err)
			return
		}

		if head.GetDataLen() > 0 {
			msg := head.(*network.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("recv ID: %d LEN: %d DATA: %s", msg.ID, msg.DataLen, msg.Data)
		}

		time.Sleep(time.Second)
	}
}
