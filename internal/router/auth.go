package router

import (
	"app/internal/handler"
	"app/internal/service"
	"app/pkg/tracker"
	"net/http"
)

func addAuthRoutes(
	mux *http.ServeMux,
	authService service.AuthService,
	tracker *tracker.Tracker,
) {
	authHandler := handler.NewAuthHandler(authService, tracker)

	mux.HandleFunc("GET /auth/password", authHandler.PasswordLogInPage)
	mux.HandleFunc("POST /auth/password", authHandler.PasswordLogIn)

	// QRcode login page
	mux.HandleFunc("GET /auth/password/qrcode", authHandler.QRcodeLogInPage)
	mux.HandleFunc("POST /auth/password/qrcode", authHandler.QRcodeLogIn)

	mux.HandleFunc("/auth/logout", authHandler.Logout)
}
