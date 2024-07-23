package routes

import (
	"app/assets"
	"app/internal/middleware"
	"app/internal/routemodule"
	"app/routes/auth"
	"app/routes/home"
	"app/routes/notfound"
	"app/routes/users"
	"net/http"
)

func Handler() http.Handler {
	// Add logging before each middleware or handler setup
	middlewareStack := middleware.CreateStack(
		middleware.Security,
		middleware.Logging,
		middleware.Authentication,
		middleware.AuthRedirect,
	)

	r := http.NewServeMux()

	// public static assets file server
	staticFS := http.FileServer(assets.Assets)
	r.Handle("/static/", staticFS)

	// add routes
	routeModules := []routemodule.RouteModule{
		users.NewUserModule("/users"),
		auth.NewAuthModule("/auth"),
	}

	for _, rm := range routeModules {
		rm.AddRoutes(r, rm.GetPrefix())
	}

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/" {
			home.Handler(w, r)
			return
		}
		notfound.Handler(w, r)
	})

	return middlewareStack(r)
}
