package middleware

import (
	"app/internal/model"
	"app/pkg/reqcontext"
	"net/http"
)

func AuthRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicRouteRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		_, ok := r.Context().Value(reqcontext.ReqContextKeyUser).(model.User)
		if !ok {
			http.Redirect(w, r, "/auth/password", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})

}
