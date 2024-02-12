package login

import (
	"github.com/gorilla/mux"
)

func AddRouter(r *mux.Router) {
	s := r.PathPrefix("/login").Subrouter()

	s.HandleFunc("/password", passwordLoginViewHandler).Methods("GET")
	s.HandleFunc("/password", passwordLoginHandler).Methods("POST")

}
