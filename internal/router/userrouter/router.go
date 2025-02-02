package userrouter

import (
	"app/internal/handlers/userhandler"
	"app/internal/services/userservice"
	"net/http"
)

func NewRouter(userService userservice.UserService) http.Handler {
	mux := http.NewServeMux()

	userHandler := userhandler.NewUserHandler(userService)

	mux.HandleFunc("GET /users", userHandler.UsersHomePage)

	mux.HandleFunc("GET /users/add", userHandler.AddUserPage)
	mux.HandleFunc("POST /users/add", userHandler.AddUser)

	mux.HandleFunc("GET /users/add-api-user", userHandler.AddAPIUserPage)
	mux.HandleFunc("POST /users/add-api-user", userHandler.AddAPIUser)

	mux.HandleFunc("GET /users/{id}", userHandler.UserPage)

	mux.HandleFunc("GET /users/{id}/edit", userHandler.EditUserPage)
	mux.HandleFunc("POST /users/{id}/edit", userHandler.EditUser)

	mux.HandleFunc("GET /users/{id}/reset-password", userHandler.ResetPasswordPage)
	mux.HandleFunc("POST /users/{id}/reset-password", userHandler.ResetPassword)

	return mux
}
