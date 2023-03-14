package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/atlasfox007/Golang-Simple-Backend-App/handler"
	"github.com/atlasfox007/Golang-Simple-Backend-App/repository"
	"github.com/atlasfox007/Golang-Simple-Backend-App/services"
)

const (
	mongoDBHost     = "localhost"
	mongoDBPort     = "27017"
	mongoDBDatabase = "mydatabase"
)

func main() {
	// create a MongoDB client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + mongoDBHost + ":" + mongoDBPort))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	// create a UserRepository implementation using the MongoDB client
	repo := repository.NewUserRepository(client.Database(mongoDBDatabase).Collection("users"))

	// create a UserService implementation using the UserRepository implementation
	service := services.NewUserService(repo)

	// create a UserHandler using the UserService implementation
	handler := handler.NewUserHandler(service)

	// create a router and register the UserHandler's routes
	router := mux.NewRouter()
	handler.RegisterRoutes(router)

	// create a server with a graceful shutdown
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		log.Println("Starting server...")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// wait for an interrupt signal to shutdown the server gracefully
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	log.Println("Shutting down server...")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Server stopped.")
}
