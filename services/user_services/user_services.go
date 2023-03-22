package user_services

import (
	"errors"
	"github.com/atlasfox007/Golang-Simple-Backend-App/repository/user_repository"

	"github.com/atlasfox007/Golang-Simple-Backend-App/model"
)

type UserService interface {
	GetAllUsers() ([]model.User, error)
	GetUserByID(id string) (*model.User, error)
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	DeleteUser(id string) error
}

type userService struct {
	repo user_repository.UserRepository
}

func NewUserService(repo user_repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) GetAllUsers() ([]model.User, error) {
	return s.repo.GetAllUsers()
}

func (s *userService) GetUserByID(id string) (*model.User, error) {
	if id == "" {
		return nil, errors.New("ID cannot be empty")
	}
	return s.repo.GetUserByID(id)
}

func (s *userService) CreateUser(user *model.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	return s.repo.CreateUser(user)
}

func (s *userService) UpdateUser(user *model.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	return s.repo.UpdateUser(user)
}

func (s *userService) DeleteUser(id string) error {
	if id == "" {
		return errors.New("ID cannot be empty")
	}
	return s.repo.DeleteUser(id)
}
