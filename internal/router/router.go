package router

import (
	"app/assets"
	"app/internal/service"
	"app/internal/views/homeview"
	"app/internal/views/notfoundview"
	"app/pkg/apphmac"
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
	FileService             service.FileService
	GalleryService          service.GalleryService
	NotificationService     service.NotificationService
	PDFService              service.PDFService
	PrintNodeService        service.PrintNodeService
	ResourceService         service.ResourceService
	SearchService           service.SearchService
	ServicesService         service.ServicesService
	StockTransactionService service.StockTransactionService
	StockItemService        service.StockItemService
	TeamService             service.TeamService
	UserService             service.UserService
}

func NewRouter(services *Services, appHMAC apphmac.AppHMAC) http.Handler {

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

	// add routes
	addAuthRoutes(mux, services.AuthService)
	addAndonRoutes(
		mux,
		services.AndonService,
		services.AndonIssueService,
		services.CommentService,
		services.GalleryService,
		services.TeamService,
		appHMAC,
	)
	addAndonIssueRoutes(mux, services.AndonIssueService, services.TeamService)
	addCameraScannerRoutes(mux)
	addFileRoutes(mux, services.FileService)
	addGalleryRoutes(mux, services.GalleryService, appHMAC)
	addNotificationRoutes(mux, services.NotificationService)
	addPDFRoutes(mux, services.PDFService, services.PrintNodeService)
	addResourceRoutes(
		mux,
		services.GalleryService,
		services.ResourceService,
		services.ServicesService,
		services.CommentService,
		services.TeamService,
		appHMAC,
	)
	addSearchRoutes(mux, services.SearchService)
	addCommentRoutes(mux, services.CommentService, services.FileService, appHMAC)
	addServiceRoutes(
		mux,
		services.GalleryService,
		services.ResourceService,
		services.ServicesService,
		services.CommentService,
		services.TeamService,
		appHMAC,
	)
	addStockItemRoutes(mux, services.StockItemService, services.CommentService, services.GalleryService, appHMAC)
	addStockTransactionRoutes(mux, services.StockItemService, services.StockTransactionService)
	addTeamRoutes(mux, services.TeamService, services.UserService)
	addUserRoutes(mux, services.UserService)

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

	})

	return middlewareStack(mux)
}
