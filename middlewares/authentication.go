package middlewares

import (
	"app/db"
	"app/src/auth"
	userModel "app/src/users/model"
	"app/utils"
	"context"
	"net/http"
	"strings"
)

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id int
		cookie, err := r.Cookie("login-session")
		if err == nil {
			err = utils.CookieInstance.Decode("login-session", cookie.Value, &id)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
		}

		if id == 0 {
			params := r.URL.Query()
			authToken := params.Get("authToken")
			if authToken == "" {
				authToken = r.Header.Get("Authorization")
			}

			if authToken == "" {
				authToken = r.Header.Get("X-Authorization")
			}

			if authToken != "" {
				apiUsername, apiPassword, found := strings.Cut(authToken, ":")
				if !found {
					next.ServeHTTP(w, r)
					return
				}
				authUser, err := auth.VerifyUser(apiUsername, apiPassword)
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}
				id = authUser.UserId
			}
		}

		// If id is still 0, then the user is not authenticated
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		db := db.UseDB()
		user, err := userModel.ByID(db, id)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), utils.ContextKeyUser, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
