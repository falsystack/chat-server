package network

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

type Network struct {
	engin *gin.Engine
}

func NewServer() *Network {
	n := &Network{engin: gin.New()}

	n.engin.Use(gin.Logger())   // user logger
	n.engin.Use(gin.Recovery()) // panic 또는 에러로 인한 서버가 죽으면 다시 기동시켜주는 역할을 한다
	n.engin.Use(cors.New(cors.Config{
		AllowWebSockets: true,
		AllowOrigins:    []string{"*"},
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:    []string{"*"},
	}))

	r := NewRoom()
	go r.RunInit()

	n.engin.GET("/room", r.SocketServe)

	return n
}

func (n *Network) StartServer() error {
	log.Println("Starting Server...")
	return n.engin.Run(":8080")
}
