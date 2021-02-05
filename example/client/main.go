package main

import (
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
		_, err := conn.Write([]byte(time.Now().String()))
		if err != nil {
			log.Println(err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("server callback : %s %d", buf, cnt)

		time.Sleep(time.Second)
	}
}
