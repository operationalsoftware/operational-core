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
	s.HandleFunc("/add/validate", validateAddUserHandler).Methods("POST")

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

	// User page
	s.HandleFunc("/{id}", userViewHandler).Methods("GET")

	// Edit user
	s.HandleFunc("/{id}/edit", editUserViewHandler).Methods("GET")
	s.HandleFunc("/{id}/edit", editUserHandler).Methods("POST")
	s.HandleFunc("/{id}/edit/validate", validateEditUserHandler).Methods("POST")
}
