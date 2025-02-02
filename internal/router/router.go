package router

import (
	"app/assets"
	"app/internal/router/authrouter"
	"app/internal/router/userrouter"
	"app/internal/services/authservice"
	"app/internal/services/userservice"
	"app/internal/views/homeview"
	"app/internal/views/notfoundview"
	"app/pkg/middleware"
	"app/pkg/reqcontext"
	"net/http"
)

type Services struct {
	AuthService authservice.AuthService
	UserService userservice.UserService
}

func NewRouter(services *Services) http.Handler {

	// create the Authentication middleware with dependency injection
	authenticationMiddleware := middleware.NewAuthenticationMiddleware(
		services.AuthService,
		services.UserService,
	)

	// Add logging before each middleware or handler setup
	middlewareStack := middleware.CreateStack(
		middleware.Security,
		middleware.Logging,
		authenticationMiddleware.Authentication,
		middleware.AuthRedirect,
	)

	mux := http.NewServeMux()

	// public static assets file server
	staticFS := http.FileServer(assets.Assets)
	mux.Handle("/static/", staticFS)

	mux.Handle("/auth/", authrouter.NewRouter(services.AuthService))
	mux.Handle("/user/", userrouter.NewRouter(services.UserService))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := reqcontext.GetContext(r)

		if r.Method == http.MethodGet && r.URL.Path == "/" {
			_ = homeview.HomePage(&homeview.HomePageProps{
				Ctx: ctx,
			}).Render(w)
			return
		}

		_ = notfoundview.NotFoundPage(&notfoundview.NotFoundPageProps{
			Ctx: ctx,
		}).Render(w)

		return
	})

	return middlewareStack(mux)
}
