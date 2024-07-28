package main

import "chat-server/network"

func main() {
	server := network.NewServer()
	server.StartServer()
}
