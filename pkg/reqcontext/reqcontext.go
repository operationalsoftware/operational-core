package reqcontext

import (
	"app/internal/models"
	"net/http"
)

type ReqContext struct {
	User models.User
	Req  *http.Request
}

type ReqContextKey string

const (
	ReqContextKeyUser ReqContextKey = "user"
)

func GetContext(r *http.Request) ReqContext {
	// Get the user from the context
	user, ok := r.Context().Value(ReqContextKeyUser).(models.User)
	if !ok {
		user = models.User{}
	}

	return ReqContext{
		User: user,
		Req:  r,
	}
}
