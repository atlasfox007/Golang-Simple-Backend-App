package user_repository

import (
	"github.com/atlasfox007/Golang-Simple-Backend-App/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	GetAllUsers() ([]model.User, error)
	GetUserByID(id string) (*model.User, error)
	CreateUser(user *model.User) error
	UpdateUser(user *model.User) error
	DeleteUser(id string) error
	GetByUsername(username string) (*model.User, error)
}

type userMongoRepository struct {
	collection *mongo.Collection
}

// Placeholder for other possible database used
type userOtherRepository struct {
}
