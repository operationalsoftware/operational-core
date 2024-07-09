package middleware

import (
	"app/models/usermodel"
	"app/reqcontext"
	"net/http"
)

func AuthRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(reqcontext.ReqContextKeyUser).(usermodel.User)
		if !ok {
			http.Redirect(w, r, "/auth/password", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})

}
