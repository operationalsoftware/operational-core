package users

import (
	"app/routes/notfound"
	"app/routes/users/usershandlers"
	"net/http"
)

func Handler() http.Handler {
	r := http.NewServeMux()

	r.HandleFunc("GET /add", usershandlers.AddUserPage)
	r.HandleFunc("POST /add", usershandlers.AddUser)

	r.HandleFunc("GET /add-api-user", usershandlers.AddAPIUserPage)
	r.HandleFunc("POST /add-api-user", usershandlers.AddAPIUser)

	r.HandleFunc("GET /{id}", usershandlers.UserPage)

	r.HandleFunc("GET /{id}/edit", usershandlers.EditUserPage)
	r.HandleFunc("POST /{id}/edit", usershandlers.EditUser)

	r.HandleFunc("GET /{id}/reset-password", usershandlers.ResetPasswordPage)
	r.HandleFunc("POST /{id}/reset-password", usershandlers.ResetPassword)

	r.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			notfound.Handler(w, r)
			return
		}
		usershandlers.UsersHomePage(w, r)
	})

	r.HandleFunc("/", notfound.Handler)

	return r
}
