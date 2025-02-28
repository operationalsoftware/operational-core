package router

import (
	"app/internal/handlers/authhandler"
	"app/internal/service"
	"net/http"
)

func addAuthRoutes(
	mux *http.ServeMux,
	authService service.AuthService,
) {
	authHandler := authhandler.NewAuthHandler(authService)

	mux.HandleFunc("GET /auth/password", authHandler.PasswordLogInPage)
	mux.HandleFunc("POST /auth/password", authHandler.PasswordLogIn)

	mux.HandleFunc("/logout", authHandler.Logout)
}
