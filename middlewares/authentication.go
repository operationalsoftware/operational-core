package middlewares

import (
	"context"
	"net/http"
	"app/db"
	userModel "app/src/users/model"
	"app/utils"
)

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id int
		if cookie, err := r.Cookie("login-session"); err != nil {
			next.ServeHTTP(w, r)
			return
		} else {
			if err = utils.CookieInstance.Decode("login-session", cookie.Value, &id); err != nil {
				next.ServeHTTP(w, r)
				return
			}
		}

		if id <= 0 {
			next.ServeHTTP(w, r)
			return
		}

		db := db.UseDB()
		user, err := userModel.ByID(db, id)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
