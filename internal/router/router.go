package router

import (
	"app/assets"
	"app/internal/service"
	"app/internal/views/camerascannerview"
	"app/internal/views/homeview"
	"app/internal/views/notfoundview"
	"app/pkg/middleware"
	"app/pkg/reqcontext"
	"fmt"
	"net/http"
)

type Services struct {
	AndonService            service.AndonService
	AndonIssueService       service.AndonIssueService
	AuthService             service.AuthService
	CommentService          service.CommentService
	UserService             service.UserService
	StockTransactionService service.StockTransactionService
	SearchService           service.SearchService
	PDFService              service.PDFService
	FileService             service.FileService
	TeamService             service.TeamService
	StockItemService        service.StockItemService
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
	addAuthRoutes(mux, services.AuthService)
	addAndonRoutes(mux, services.AndonService, services.AndonIssueService, services.CommentService, services.TeamService, services.FileService)
	addAndonIssueRoutes(mux, services.AndonIssueService, services.TeamService)
	addCommentRoutes(mux, services.CommentService)
	addFileRoutes(mux, services.FileService)
	addPDFRoutes(mux, services.PDFService)
	addSearchRoutes(mux, services.SearchService)
	addTeamRoutes(mux, services.TeamService, services.UserService)
	addStockTransactionRoutes(mux, services.StockItemService, services.StockTransactionService)
	addStockItemRoutes(mux, services.StockItemService, services.CommentService, services.FileService)
	addUserRoutes(mux, services.UserService)

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
