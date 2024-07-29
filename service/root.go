package service

import (
	"chat-server/repository"
	"chat-server/types/schema"
	"log"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repository *repository.Repository) *Service {
	s := &Service{repository: repository}

	return s
}

func (s *Service) EnterRoom(roomName string) ([]*schema.Chat, error) {
	if res, err := s.repository.GetChatList(roomName); err != nil {
		log.Println("Failed to get all chat list", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}

func (s *Service) RoomList() ([]*schema.Room, error) {
	if res, err := s.repository.RoomList(); err != nil {
		log.Println("Failed to get all room list", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}

func (s *Service) MakeRoom(name string) error {
	if err := s.repository.MakeRoom(name); err != nil {
		log.Println("Failed to make new room", "err", err.Error())
		return err
	} else {
		return nil
	}
}

func (s *Service) Room(name string) (*schema.Room, error) {
	if res, err := s.repository.Room(name); err != nil {
		log.Println("Failed to get room", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}
