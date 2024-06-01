package main

import "chat-server/network"

func main() {
	n := network.NewServer()
	n.StartServer()
}
