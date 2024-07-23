package users

import (
	"app/internal/routemodule"
	"app/routes/users/usershandlers"
	"net/http"
)

type UserModule struct {
	routemodule.PrefixedRouteModule
}

func NewUserModule(prefix string) *UserModule {
	return &UserModule{PrefixedRouteModule: routemodule.PrefixedRouteModule{Prefix: prefix}}
}

func (u *UserModule) AddRoutes(r *http.ServeMux, prefix string) {
	u.Prefix = prefix

	r.HandleFunc("GET "+prefix, usershandlers.UsersHomePage)

	r.HandleFunc("GET "+prefix+"/add", usershandlers.AddUserPage)
	r.HandleFunc("POST "+prefix+"/add", usershandlers.AddUser)

	r.HandleFunc("GET "+prefix+"/add-api-user", usershandlers.AddAPIUserPage)
	r.HandleFunc("POST "+prefix+"/add-api-user", usershandlers.AddAPIUser)

	r.HandleFunc("GET "+prefix+"/{id}", usershandlers.UserPage)

	r.HandleFunc("GET "+prefix+"/{id}/edit", usershandlers.EditUserPage)
	r.HandleFunc("POST "+prefix+"/{id}/edit", usershandlers.EditUser)

	r.HandleFunc("GET "+prefix+"/{id}/reset-password", usershandlers.ResetPasswordPage)
	r.HandleFunc("POST "+prefix+"/{id}/reset-password", usershandlers.ResetPassword)
}
