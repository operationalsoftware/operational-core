package middleware

import (
	"app/internal/reqcontext"
	"app/models/usermodel"
	"net/http"
)

func AuthRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicRouteRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		_, ok := r.Context().Value(reqcontext.ReqContextKeyUser).(usermodel.User)
		if !ok {
			http.Redirect(w, r, "/auth/password", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})

}
