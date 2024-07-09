package middleware

import (
	"app/db"
	"app/models/authmodel"
	"app/models/usermodel"
	"app/reqcontext"
	"app/utils"
	"context"
	"database/sql"
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
