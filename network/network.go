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
	port       string
	ip         string
}

// NewServer Constructor
func NewServer(service *service.Service, rep *repository.Repository, port string) *Server {
	s := &Server{
		engine:     gin.New(),
		service:    service,
		repository: rep,
		port:       port,
	}
	s.engine.Use(gin.Logger())
	s.engine.Use(gin.Recovery()) // panic 등으로 인해 서버가 종료되었을 때 재시작시켜준다.
	s.engine.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"*"},
		AllowHeaders:    []string{"*"},
		AllowWebSockets: true,
	}))

	registerServer(s.engine)

	return s
}

func (s *Server) StartServer() error {
	log.Println("Starting server...")
	return s.engine.Run(s.port)
}
