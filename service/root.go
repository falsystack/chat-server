package service

import (
	"chat-server/repository"
	"chat-server/types/schema"
	"log"
)

type Service struct {
	rep *repository.Repository
}

func NewService(rep *repository.Repository) *Service {
	s := &Service{rep: rep}

	return s
}

func (s *Service) EnterRoom(roomName string) ([]*schema.Chat, error) {
	if res, err := s.rep.GetChatList(roomName); err != nil {
		log.Println("Failed To Get Chat List", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}

func (s *Service) RoomList() ([]*schema.Room, error) {
	if res, err := s.RoomList(); err != nil {
		log.Println("Failed To Get All Room  List", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}

func (s *Service) MakeRoom(name string) error {
	if err := s.rep.MakeRoom(name); err != nil {
		log.Println("Failed To Make Room", "err", err.Error())
		return err
	} else {
		return nil
	}
}

func (s *Service) Room(name string) (*schema.Room, error) {
	if res, err := s.Room(name); err != nil {
		log.Println("Failed To Get Room", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}
