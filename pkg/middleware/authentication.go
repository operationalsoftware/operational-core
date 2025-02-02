package middleware

import (
	"app/internal/models"
	"app/internal/services/authservice"
	"app/internal/services/userservice"
	"app/pkg/cookie"
	"app/pkg/reqcontext"
	"context"
	"net/http"
	"strings"
)

type AuthenticationMiddleware struct {
	authService authservice.AuthService
	userService userservice.UserService
}

func NewAuthenticationMiddleware(
	authService authservice.AuthService,
	userService userservice.UserService,
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

				var out models.VerifyPasswordLoginOutput

				out, err = m.authService.VerifyPasswordLogin(
					r.Context(),
					models.VerifyPasswordLoginInput{
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
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), reqcontext.ReqContextKeyUser, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
