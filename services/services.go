package services

import (
	"github.com/atlasfox007/Golang-Simple-Backend-App/repository"
	"github.com/atlasfox007/Golang-Simple-Backend-App/services/user_auth"
	"github.com/atlasfox007/Golang-Simple-Backend-App/services/user_services"
)

type Services struct {
	UserServices user_services.UserService
	UserAuth     user_auth.UserAuth
}

func InitServices(repositories *repository.Repositories) *Services {
	return &Services{
		UserServices: user_services.NewUserService(repositories.UserRepository),
		UserAuth:     user_auth.NewUserAuthService(repositories.UserRepository),
	}
}
