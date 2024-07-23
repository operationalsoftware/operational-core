package middleware

import (
	"app/internal/cookie"
	"app/internal/db"
	"app/internal/reqcontext"
	"app/models/authmodel"
	"app/models/usermodel"
	"context"
	"database/sql"
	"net/http"
	"strings"
)

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isPublicRouteRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		var id int
		sesscookie, err := r.Cookie("login-session")
		if err == nil {
			err = cookie.CookieInstance.Decode("login-session", sesscookie.Value, &id)
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
				ex := db.UseDB()
				var out authmodel.VerifyPasswordLoginOutput

				err = db.WithTx(ex, func(tx *sql.Tx) error {
					out, err = authmodel.VerifyPasswordLogin(tx, authmodel.VerifyPasswordLoginInput{
						Username: apiUsername,
						Password: apiPassword,
					})
					return err
				})

				if err != nil {
					next.ServeHTTP(w, r)
					return
				}
				id = out.AuthUser.UserId
			}
		}

		// If id is still 0, then the user is not authenticated
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		db := db.UseDB()
		user, err := usermodel.ByID(db, id)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), reqcontext.ReqContextKeyUser, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
