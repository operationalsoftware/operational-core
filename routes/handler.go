package routes

import (
	"app/assets"
	"app/middleware"
	"app/routes/auth"
	"app/routes/home"
	"app/routes/notfound"
	"app/routes/users"
	"net/http"
)

func Handler() http.Handler {

	// all routes middlewares
	pubstack := middleware.CreateStack(
		middleware.Security,
		middleware.Logging,
	)

	// private routes middleware
	privstack := middleware.CreateStack(
		middleware.Authentication,
		middleware.AuthRedirect,
	)

	// main router
	r := http.NewServeMux()

	// public static assets file server
	staticFS := http.FileServer(assets.Assets)
	r.Handle("/static/", staticFS)

	// public auth routes (log in, set/reset password)
	r.Handle("/auth/", http.StripPrefix("/auth", auth.Handler()))

	// private routes
	pr := http.NewServeMux()

	pr.Handle("/users/", http.StripPrefix("/users", users.Handler()))

	pr.HandleFunc("GET /", home.HomePage)

	pr.HandleFunc("/", notfound.Handler)

	// use private routes handler for all remaining routes
	r.Handle("/", privstack(pr))

	return pubstack(r)
}
