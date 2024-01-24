package router

import (
	"operationalcore/handlers"

	"github.com/gorilla/mux"
)

func AddUserRouter(r *mux.Router) {
	s := r.PathPrefix("/users").Subrouter()

	// Users home
	s.HandleFunc("", handlers.UsersPage).Methods("GET")

	// Add user
	s.HandleFunc("/add", handlers.CreateUserPage).Methods("GET")
	s.HandleFunc("/add", handlers.CreateUser).Methods("POST")

	// Create user form
	s.HandleFunc("/add/first-name", handlers.CreateUserFormFirstName).Methods("POST")
	s.HandleFunc("/add/last-name", handlers.CreateUserFormLastName).Methods("POST")
	s.HandleFunc("/add/email", handlers.CreateUserFormEmail).Methods("POST")
	s.HandleFunc("/add/username", handlers.CreateUserFormUsername).Methods("POST")
	s.HandleFunc("/add/password", handlers.CreateUserFormPassword).Methods("POST")
	s.HandleFunc("/add/confirm-password", handlers.CreateUserConfirmPassword).Methods("POST")

	// User page
	s.HandleFunc("/{id}", handlers.UserPage).Methods("GET")

	// Edit user
	s.HandleFunc("/{id}/edit", handlers.EditUserPage).Methods("GET")
	s.HandleFunc("/{id}/edit", handlers.EditUser).Methods("POST")
}
