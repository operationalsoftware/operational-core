package auth

import (
	"app/routes/auth/authhandlers"
	"app/routes/notfound"
	"fmt"
	"net/http"
)

func Handler() http.Handler {
	r := http.NewServeMux()

	// log in with password
	r.HandleFunc("GET /password", authhandlers.PasswordLogInPage)
	r.HandleFunc("POST /password", authhandlers.PasswordLogIn)

	// log out
	r.HandleFunc("/logout", authhandlers.Logout)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("NOT FOUND")
		notfound.Handler(w, r)
	})

	return r
}
