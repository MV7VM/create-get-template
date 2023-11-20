package service

import (
	"strconv"
)

type repository interface {
	CreateTemplate(msg string, id int) (string, error)
	GetTemplate(id int) (string, error)
	GetTemplates() (string, error)
}

type Service struct {
	repository
}

func New(repository repository) *Service {
	return &Service{repository}
}

func (s *Service) Create(msg string, prefid string) (string, error) {
	if Id, ok := strconv.Atoi(prefid); ok != nil {
		//fmt.Println("Validate - ", msg, Id)
		return s.repository.CreateTemplate(msg, -1)
	} else {
		return s.repository.CreateTemplate(msg, Id)
	}
}

func (s *Service) Get(id string) (string, error) {
	Id, err := strconv.Atoi(id)
	if err != nil {
		return "", err
	}
	return s.repository.GetTemplate(Id)
}

func (s *Service) GetAll() (string, error) {
	return s.repository.GetTemplates()
}
