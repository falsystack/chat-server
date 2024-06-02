package service

import "chat-server/repository"

type Service struct {
	rep *repository.Repository
}

func NewService(rep *repository.Repository) *Service {
	s := &Service{rep: rep}

	return s
}
