package model

type model struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Service interface {
	Add(name string) error
	Remove(id int) error
	GetAll() ([]model, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) Add(name string) error {
	return nil
}

func (s *service) Remove(id int) error {
	return nil
}

func (s *service) GetAll() ([]model, error) {
	return []model{}, nil
}
