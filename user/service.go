package user

type repository interface {
	CreateUser(name string) (*User, error)
	UpdateUser(u *User) error
	GetUser(id int) (*User, error)
}

type Service struct {
	repository repository
}

func NewService(repository repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateUser(name string) (*User, error) {
	return s.repository.CreateUser(name)
}

func (s *Service) UpdateUser(u *User) error {
	return s.repository.UpdateUser(u)
}

func (s *Service) GetUser(id int) (*User, error) {
	return s.repository.GetUser(id)
}
