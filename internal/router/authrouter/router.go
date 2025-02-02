package authrouter

import (
	"app/internal/handlers/authhandler"
	"app/internal/services/authservice"
	"net/http"
)

func NewRouter(authService authservice.AuthService) http.Handler {
	mux := http.NewServeMux()

	authHandler := authhandler.NewAuthHandler(authService)

	mux.HandleFunc("GET /auth/password", authHandler.PasswordLogInPage)
	mux.HandleFunc("POST /auth/password", authHandler.PasswordLogIn)

	// log out (any method)
	mux.HandleFunc("/logout", authHandler.Logout)

	return mux
}
