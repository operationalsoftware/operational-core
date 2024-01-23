package middlewares

import (
	"context"
	"net/http"
	"operationalcore/db"
	"operationalcore/model"
	"operationalcore/utils"
	"strconv"
)

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var id int
		if cookie, err := r.Cookie("operationalcore-session"); err != nil {
			next.ServeHTTP(w, r)
			return
		} else {
			if err = utils.CookieInstance.Decode("operationalcore-session", cookie.Value, &id); err != nil {
				next.ServeHTTP(w, r)
				return
			}
		}

		if id <= 0 {
			next.ServeHTTP(w, r)
			return
		}

		dbInstance := db.UseDB()
		user := model.GetUser(dbInstance, strconv.Itoa(id))

		if user.UserId == 0 {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
