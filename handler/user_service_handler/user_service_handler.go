package user_service_handler

import (
	"encoding/json"
	"github.com/atlasfox007/Golang-Simple-Backend-App/services/user_services"
	"net/http"

	"github.com/atlasfox007/Golang-Simple-Backend-App/middleware"
	"github.com/atlasfox007/Golang-Simple-Backend-App/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserServiceHandler struct {
	service user_services.UserService
}

func NewUserServiceHandler(service user_services.UserService) *UserServiceHandler {
	return &UserServiceHandler{
		service: service,
	}
}

func (h *UserServiceHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users", middleware.IsAuthenticated(h.getAllUsersHandler)).Methods("GET")
	r.HandleFunc("/users/{id}", middleware.IsAuthenticated(h.getUserByIDHandler)).Methods("GET")
	r.HandleFunc("/users/{id}", middleware.IsAuthenticated(h.updateUserHandler)).Methods("PUT")
	r.HandleFunc("/users/{id}", middleware.IsAuthenticated(h.deleteUserHandler)).Methods("DELETE")
}

func (h *UserServiceHandler) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserServiceHandler) getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.service.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserServiceHandler) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user model.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = h.service.UpdateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserServiceHandler) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := h.service.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
