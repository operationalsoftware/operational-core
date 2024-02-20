package utils

import (
	userModel "app/src/users/model"
	"net/http"
)

type Context struct {
	User userModel.User
	Req  *http.Request
}

type ContextKey string

const (
	ContextKeyUser ContextKey = "user"
)

func GetContext(r *http.Request) Context {
	// Get the user from the context
	user, ok := r.Context().Value(ContextKeyUser).(userModel.User)
	if !ok {
		user = userModel.User{}
	}

	return Context{
		User: user,
		Req:  r,
	}
}
