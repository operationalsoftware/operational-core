package router

import (
	"app/assets"
	"app/internal/service"
	"app/internal/views/camerascannerview"
	"app/internal/views/homeview"
	"app/internal/views/notfoundview"
	"app/pkg/middleware"
	"app/pkg/reqcontext"
	"app/pkg/tracker"
	"fmt"
	"net/http"
)

type Services struct {
	AuthService             service.AuthService
	UserService             service.UserService
	StockTransactionService service.StockTransactionService
	SearchService           service.SearchService
	PDFService              service.PDFService
	FileService             service.FileService
	Tracker                 *tracker.Tracker
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

	// add routes for services
	addAuthRoutes(mux, services.AuthService, services.Tracker)
	addUserRoutes(mux, services.UserService)
	addStockTransactionRoutes(mux, services.StockTransactionService)
	addSearchRoutes(mux, services.SearchService)
	addPDFRoutes(mux, services.PDFService)
	addFileRoutes(mux, services.PDFService)

	// Camera scanner route
	mux.HandleFunc("/camera-scanner", func(w http.ResponseWriter, r *http.Request) {
		ctx := reqcontext.GetContext(r)

		_ = camerascannerview.CameraScannerApp(&camerascannerview.CameraScannerAppProps{
			Ctx: ctx,
		}).Render(w)

	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := reqcontext.GetContext(r)

		// get homepage
		if r.Method == http.MethodGet && r.URL.Path == "/" {
			_ = homeview.HomePage(&homeview.HomePageProps{
				Ctx: ctx,
			}).Render(w)
			return
		}

		// catch all - Not Found
		w.WriteHeader(http.StatusNotFound)

		if r.Method == http.MethodGet {
			_ = notfoundview.NotFoundPage(&notfoundview.NotFoundPageProps{
				Ctx: ctx,
			}).Render(w)
			return
		}

		fmt.Fprintln(w, "404 Not Found")

		return
	})

	return middlewareStack(mux)
}
