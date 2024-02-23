package reqContext

import (
	userModel "app/src/users/model"
	"net/http"
)

type ReqContext struct {
	User userModel.User
	Req  *http.Request
}

type ReqContextKey string

const (
	ReqContextKeyUser ReqContextKey = "user"
)

func GetContext(r *http.Request) ReqContext {
	// Get the user from the context
	user, ok := r.Context().Value(ReqContextKeyUser).(userModel.User)
	if !ok {
		user = userModel.User{}
	}

	return ReqContext{
		User: user,
		Req:  r,
	}
}
