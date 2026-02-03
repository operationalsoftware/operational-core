package router

import (
	"app/internal/handler"
	"app/internal/service"
	"net/http"
)

func addAuthRoutes(
	mux *http.ServeMux,
	authService service.AuthService,
	notificationService service.NotificationService,
) {
	authHandler := handler.NewAuthHandler(authService, notificationService)

	mux.HandleFunc("GET /auth/password", authHandler.PasswordLogInPage)
	mux.HandleFunc("POST /auth/password", authHandler.PasswordLogIn)

	// QRcode login page
	mux.HandleFunc("GET /auth/password/qrcode", authHandler.QRcodeLogInPage)
	mux.HandleFunc("POST /auth/password/qrcode", authHandler.QRcodeLogIn)

	mux.HandleFunc("/auth/logout", authHandler.Logout)

	mux.HandleFunc("GET /auth/microsoft/login", authHandler.MicrosoftLogin)
	mux.HandleFunc("GET /auth/microsoft/callback", authHandler.MicrosoftCallback)
}
