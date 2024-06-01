package network

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

type Network struct {
	engine *gin.Engine
}

// NewServer Constructor
func NewServer() *Network {
	n := &Network{
		engine: gin.New(),
	}
	n.engine.Use(gin.Logger())
	n.engine.Use(gin.Recovery()) // panic 등으로 인해 서버가 종료되었을 때 재시작시켜준다.
	n.engine.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"*"},
		AllowHeaders:    []string{"*"},
		AllowWebSockets: true,
	}))

	return n
}

func (n *Network) StartServer() error {
	log.Println("Starting server...")
	return n.engine.Run(":8080")
}
