package reqcontext

import (
	"app/internal/model"
	"net/http"
)

type ReqContext struct {
	User model.User
	Req  *http.Request
}

type ReqContextKey string

const (
	ReqContextKeyUser ReqContextKey = "user"
)

func GetContext(r *http.Request) ReqContext {
	// Get the user from the context
	user, ok := r.Context().Value(ReqContextKeyUser).(model.User)
	if !ok {
		user = model.User{}
	}

	return ReqContext{
		User: user,
		Req:  r,
	}
}
