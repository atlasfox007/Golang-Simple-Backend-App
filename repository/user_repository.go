package repository

import (
	"context"
	"time"

	"github.com/atlasfox007/Golang-Simple-Backend-App/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) UserRepository {
	return &userRepository{
		collection: collection,
	}
}

func (r *userRepository) GetAllUsers() ([]model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cur, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var users []model.User
	err = cur.All(ctx, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) GetUserByID(id string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user model.User
	err = r.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"name": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *userRepository) UpdateUser(user *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.UpdatedAt = time.Now()

	oid, err := primitive.ObjectIDFromHex(user.ID.Hex())
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": user})
	return err
}

func (r *userRepository) DeleteUser(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}
