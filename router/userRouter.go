package router

import (
	"operationalcore/handlers"

	"github.com/gorilla/mux"
)

func AddUserRouter(r *mux.Router) {
	s := r.PathPrefix("/users").Subrouter()

	s.HandleFunc("", handlers.UsersPage).Methods("GET")

	s.HandleFunc("/view/{id}", handlers.UserPage).Methods("GET")

	s.HandleFunc("/create", handlers.CreateUserPage).Methods("GET")
	s.HandleFunc("/create", handlers.CreateUser).Methods("POST")

	s.HandleFunc("/create/first-name", handlers.CreateUserFormFirstName).Methods("POST")
	s.HandleFunc("/create/last-name", handlers.CreateUserFormLastName).Methods("POST")
	s.HandleFunc("/create/email", handlers.CreateUserFormEmail).Methods("POST")
	s.HandleFunc("/create/username", handlers.CreateUserFormUsername).Methods("POST")
	s.HandleFunc("/create/password", handlers.CreateUserFormPassword).Methods("POST")
	s.HandleFunc("/create/confirm-password", handlers.CreateUserConfirmPassword).Methods("POST")

	s.HandleFunc("/edit/{id}", handlers.EditUserPage).Methods("GET")
	s.HandleFunc("/edit/{id}", handlers.EditUser).Methods("POST")

}
