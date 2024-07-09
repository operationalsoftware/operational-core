package auth

import (
	"app/routes/auth/authhandlers"
	"app/routes/notfound"
	"net/http"
)

func Handler() http.Handler {
	r := http.NewServeMux()

	// log in with password
	r.HandleFunc("GET /password", authhandlers.PasswordLogInPage)
	r.HandleFunc("POST /password", authhandlers.PasswordLogIn)

	// log out
	r.HandleFunc("/logout", authhandlers.Logout)

	r.HandleFunc("/", notfound.Handler)

	return r
}
