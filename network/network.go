package network

import (
	"chat-server/repository"
	"chat-server/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

type Server struct {
	engine *gin.Engine

	service    *service.Service
	repository *repository.Repository

	port string
	ip   string
}

func NewServer(service *service.Service, repository *repository.Repository, port string) *Server {
	s := &Server{
		engine:     gin.New(),
		service:    service,
		repository: repository,
		port:       port,
	}

	s.engine.Use(gin.Logger())   // user logger
	s.engine.Use(gin.Recovery()) // panic 또는 에러로 인한 서버가 죽으면 다시 기동시켜주는 역할을 한다
	s.engine.Use(cors.New(cors.Config{
		AllowWebSockets: true,
		AllowOrigins:    []string{"*"},
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:    []string{"*"},
	}))

	registerServer(s)

	return s
}

func (s *Server) StartServer() error {
	log.Println("Starting Server...")
	return s.engine.Run(s.port)
}
