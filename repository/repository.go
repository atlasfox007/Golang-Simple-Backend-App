package repository

import (
	"github.com/atlasfox007/Golang-Simple-Backend-App/repository/user_repository"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repositories struct {
	UserRepository user_repository.UserRepository
}

func InitRepositories(mongoCollection *mongo.Collection) *Repositories {
	return &Repositories{
		UserRepository: user_repository.NewUserMongoRepository(mongoCollection),
	}
}
