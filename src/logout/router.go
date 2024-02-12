package logout

import "github.com/gorilla/mux"

func AddRouter(r *mux.Router) {
	s := r.PathPrefix("/logout").Subrouter()

	s.HandleFunc("", logoutHandler).Methods("POST")
}
