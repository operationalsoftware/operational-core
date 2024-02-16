package users

import (
	"app/partials"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func AddRouter(r *mux.Router) {
	s := r.PathPrefix("/users").Subrouter()

	// Users home
	s.HandleFunc("", indexViewHandler).Methods("GET")

	// Add user
	s.HandleFunc("/add", addUserViewHandler).Methods("GET")
	s.HandleFunc("/add", addUserHandler).Methods("POST")

	// partial table
	s.HandleFunc("/table", func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		sort := values.Get("sort")
		fmt.Println("sort", sort)
		// if sort == "" {
		// 	w.Header().Set("hx-push-url", "/users/table?sort=Username-asc")
		// }
		_ = partials.UsersTable().Render(w)
	}).Methods("GET")

	// User form partials
	s.HandleFunc("/validate/first-name", validateFirstNameHandler).Methods("POST")
	s.HandleFunc("/validate/last-name", validateLastNameHandler).Methods("POST")
	s.HandleFunc("/validate/email", validateEmailHandler).Methods("POST")
	s.HandleFunc("/validate/username", validateUsernameHandler).Methods("POST")
	s.HandleFunc("/validate/password", validatePasswordHandler).Methods("POST")
	s.HandleFunc("/validate/confirm-password", validateConfirmPasswordHandler).Methods("POST")

	// User page
	s.HandleFunc("/{id}", userViewHandler).Methods("GET")

	// Reset password
	s.HandleFunc("/{id}/reset-password", resetPasswordUserViewHandler).Methods("GET")
	s.HandleFunc("/{id}/reset-password", resetPasswordUserHandler).Methods("POST")

	// Edit user
	s.HandleFunc("/{id}/edit", editUserViewHandler).Methods("GET")
	s.HandleFunc("/{id}/edit", editUserHandler).Methods("POST")
}
