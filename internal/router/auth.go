package router

import (
	"app/internal/handler"
	"app/internal/service"
	"net/http"
)

func addAuthRoutes(
	mux *http.ServeMux,
	authService service.AuthService,
) {
	authHandler := handler.NewAuthHandler(authService)

	mux.HandleFunc("GET /auth/password", authHandler.PasswordLogInPage)
	mux.HandleFunc("POST /auth/password", authHandler.PasswordLogIn)

	// QRcode login page
	mux.HandleFunc("GET /auth/password/qrcode", authHandler.QRcodeLogInPage)
	mux.HandleFunc("POST /auth/password/qrcode", authHandler.QRcodeLogIn)

	mux.HandleFunc("/auth/logout", authHandler.Logout)
}
