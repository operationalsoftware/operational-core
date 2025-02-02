package middleware

import (
	"app/internal/models"
	"app/pkg/reqcontext"
	"net/http"
)

func AuthRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicRouteRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		_, ok := r.Context().Value(reqcontext.ReqContextKeyUser).(models.User)
		if !ok {
			http.Redirect(w, r, "/auth/password", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})

}
