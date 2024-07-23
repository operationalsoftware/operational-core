package reqcontext

import (
	"app/models/usermodel"
	"net/http"
)

type ReqContext struct {
	User usermodel.User
	Req  *http.Request
}

type ReqContextKey string

const (
	ReqContextKeyUser ReqContextKey = "user"
)

func GetContext(r *http.Request) ReqContext {
	// Get the user from the context
	user, ok := r.Context().Value(ReqContextKeyUser).(usermodel.User)
	if !ok {
		user = usermodel.User{}
	}

	return ReqContext{
		User: user,
		Req:  r,
	}
}
