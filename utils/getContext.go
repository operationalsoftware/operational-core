package utils

import (
	"net/http"
	"operationalcore/model"
)

type Context struct {
	User model.User
}

func GetContext(r *http.Request) Context {
	// Get the user from the context
	user, ok := r.Context().Value("user").(model.User)
	if !ok {
		user = model.User{}
	}

	return Context{
		User: user,
	}
}
