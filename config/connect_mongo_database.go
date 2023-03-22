package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName   = "dagangan_test_db"
	mongoUri = "mongodb+srv://dagangan_test_user:dagangan12345@cluster0.tbrzgbg.mongodb.net/?retryWrites=true&w=majority"
)

var db *mongo.Database

func InitDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoUri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	db = client.Database(dbName)

	fmt.Println("Successfully connected to database.")

	return nil
}

func GetDB() *mongo.Database {
	return db
}
