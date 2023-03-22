package main

import (
	"context"
	"github.com/atlasfox007/Golang-Simple-Backend-App/handler/user_auth_handler"
	"github.com/atlasfox007/Golang-Simple-Backend-App/handler/user_service_handler"
	"github.com/atlasfox007/Golang-Simple-Backend-App/repository"
	"github.com/atlasfox007/Golang-Simple-Backend-App/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"

	"github.com/atlasfox007/Golang-Simple-Backend-App/config"
)

func main() {
	// Create database connection
	err := config.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	db := config.GetDB()
	// create a UserRepository implementation using the MongoDB client
	//repo := user_repository.NewUserMongoRepository(db.Collection("dagangan_collection"))
	repo := repository.InitRepositories(db.Collection("dagangan_collection"))

	// create a UserService implementation using the UserRepository implementation
	service := services.InitServices(repo)

	// create a UserServiceHandler using the UserService implementation
	userServiceHandler := user_service_handler.NewUserServiceHandler(service.UserServices)
	userAuthHandler := user_auth_handler.NewUserAuthHandler(service.UserAuth)

	// create a router and register the UserServiceHandler's and UserAuthHandler routes
	router := mux.NewRouter()
	userServiceHandler.RegisterRoutes(router)
	userAuthHandler.RegisterRoutes(router)

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

	// wait for an interrupt signal to shut down the server gracefully
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Server stopped.")
}
