package user_auth_handler

import (
	"encoding/json"
	"github.com/atlasfox007/Golang-Simple-Backend-App/model"
	"github.com/atlasfox007/Golang-Simple-Backend-App/services/user_auth"
	"github.com/gorilla/mux"
	"net/http"
)

type UserAuthHandler struct {
	auth user_auth.UserAuth
}

func NewUserAuthHandler(auth user_auth.UserAuth) *UserAuthHandler {
	return &UserAuthHandler{
		auth: auth,
	}
}

func (a *UserAuthHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/register", a.registerUserHandler).Methods("POST")
	r.HandleFunc("/users/login", a.loginUserHandler).Methods("POST")
}

func (a *UserAuthHandler) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = a.auth.Register(user.Name, user.Password, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	message := map[string]string{
		"success_message": "Users created successfully",
	}

	json.NewEncoder(w).Encode(message)
}

func (a *UserAuthHandler) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"name"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Authenticate the user
	tokenString, err := a.auth.Login(credentials.Username, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Return the jwt token to be used
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	message := map[string]string{
		"token": tokenString,
	}

	json.NewEncoder(w).Encode(message)
}
