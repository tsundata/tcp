package main

import "github.com/tsundata/tcp/network"

func main() {
	s := network.NewServer("example")
	s.Serve()
}
