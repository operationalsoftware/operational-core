package router

import (
	"operationalcore/handlers"

	"github.com/gorilla/mux"
)

func AddUserRouter(r *mux.Router) {
	s := r.PathPrefix("/users").Subrouter()

	s.HandleFunc("", handlers.UsersHome).Methods("GET")

	s.HandleFunc("/create", handlers.CreateUserPage).Methods("GET")
	s.HandleFunc("/create", handlers.CreateUser).Methods("POST")
	s.HandleFunc("/create/first-name", handlers.CreateUserFormFirstName).Methods("POST")
	s.HandleFunc("/create/last-name", handlers.CreateUserFormLastName).Methods("POST")
	s.HandleFunc("/create/email", handlers.CreateUserFormEmail).Methods("POST")
	s.HandleFunc("/create/username", handlers.CreateUserFormUsername).Methods("POST")

}
