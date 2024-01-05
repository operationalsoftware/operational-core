package middlewares

import (
	"net/http"
	"operationalcore/model"
	"strings"
)

func AuthRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		isLoginRoute := strings.HasPrefix(url, "/login")
		isStaticRoute := strings.HasPrefix(url, "/static")
		_, ok := r.Context().Value("user").(model.User)

		if ok && isLoginRoute {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if !ok && !isLoginRoute && !isStaticRoute {
			http.Redirect(w, r, "/login/password", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
