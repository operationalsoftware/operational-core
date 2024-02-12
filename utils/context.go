package utils

import (
	"net/http"
	userModel "app/src/users/model"
)

type Context struct {
	User userModel.User
	Req  *http.Request
}

func GetContext(r *http.Request) Context {
	// Get the user from the context
	user, ok := r.Context().Value("user").(userModel.User)
	if !ok {
		user = userModel.User{}
	}

	return Context{
		User: user,
		Req:  r,
	}
}
