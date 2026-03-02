package middleware

import (
	"app/internal/model"
	"app/pkg/reqcontext"
	"net/http"
	"strings"
)

func AuthRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicRouteRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		_, ok := r.Context().Value(reqcontext.ReqContextKeyUser).(model.User)
		if !ok {
			if strings.HasPrefix(r.URL.Path, "/api/") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/auth/password", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})

}
