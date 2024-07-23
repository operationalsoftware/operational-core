package auth

import (
	"app/internal/routemodule"
	"app/routes/auth/authhandlers"
	"net/http"
)

type AuthModule struct {
	routemodule.PrefixedRouteModule
}

func NewAuthModule(prefix string) *AuthModule {
	return &AuthModule{PrefixedRouteModule: routemodule.PrefixedRouteModule{Prefix: prefix}}
}

func (u *AuthModule) AddRoutes(r *http.ServeMux, prefix string) {
	u.Prefix = prefix

	// log in with password
	r.HandleFunc("GET "+prefix+"/password", authhandlers.PasswordLogInPage)
	r.HandleFunc("POST "+prefix+"/password", authhandlers.PasswordLogIn)

	// log out
	r.HandleFunc(prefix+"/logout", authhandlers.Logout)
}
