package main

import (
	"chat-server/config"
	"chat-server/network"
	"chat-server/repository"
	"chat-server/service"
	"flag"
	"fmt"
)

var pathFlag = flag.String("config", "./config.toml", "config set")
var port = flag.String("port", ":1010", "port set")

func main() {
	flag.Parse()
	c := config.NewConfig(*pathFlag)

	if rep, err := repository.NewRepository(c); err != nil {
		panic(err)
	} else {
		c := network.NewServer(service.NewService(rep), *port)
		c.StartServer()
	}

	fmt.Println(*pathFlag, *port)

}
