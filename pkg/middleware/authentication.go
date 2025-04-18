package middleware

import (
	"app/internal/model"
	"app/internal/service"
	"app/pkg/cookie"
	"app/pkg/reqcontext"
	"context"
	"log"
	"net/http"
	"strings"
)

type AuthenticationMiddleware struct {
	authService service.AuthService
	userService service.UserService
}

func NewAuthenticationMiddleware(
	authService service.AuthService,
	userService service.UserService,
) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		authService: authService,
		userService: userService,
	}
}

func (m *AuthenticationMiddleware) Authentication(next http.Handler) http.Handler {
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

				var out model.VerifyPasswordLoginOutput

				out, err = m.authService.VerifyPasswordLogin(
					r.Context(),
					model.VerifyPasswordLoginInput{
						Username: apiUsername,
						Password: apiPassword,
					})

				if err != nil {
					next.ServeHTTP(w, r)
					return
				}

				id = out.AuthUser.UserID
			}
		}

		// If id is still 0, then the user is not authenticated
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		user, err := m.userService.GetUserByID(r.Context(), id)
		if err != nil {
			log.Println(err)
			next.ServeHTTP(w, r)
			return
		}

		if user == nil {
			log.Println("user wth id", id, "not found")
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), reqcontext.ReqContextKeyUser, *user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
