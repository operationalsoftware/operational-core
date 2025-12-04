package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"app/internal/migrate"
	"app/internal/repository"
	"app/internal/router"
	"app/internal/service"
	"app/pkg/apphmac"
	"app/pkg/cookie"
	"app/pkg/db"
	"app/pkg/env"
	"app/pkg/filestore"
	"app/pkg/localip"
	"app/pkg/pdf"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Verify environment variables
	err := env.Verify()
	if err != nil {
		log.Fatalf("Error verifying environment: %v\n", err)
	}

	// Create database connection (SQLite - to be removed)
	// err = db.ConnectDB()
	// if err != nil {
	// 	log.Fatalf("Error connecting to SQLite: %v\n", err)
	// }
	// defer db.UseDB().Close()

	err = migrate.Run()
	if err != nil {
		log.Fatalf("fatal migration error: %v", err)
	}

	pgEnv := db.LoadPostgresEnv()
	targetConnStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		pgEnv.User, pgEnv.Password, pgEnv.Host, pgEnv.Port, pgEnv.Database)
	pgPool, err := pgxpool.New(context.Background(), targetConnStr)
	if err != nil {
		log.Fatalf("Unable to create Postgres connection pool: %v\n", err)
	}
	defer pgPool.Close() // Always close the pool when done

	// Initialise some things for start up
	err = cookie.InitCookieInstance()
	if err != nil {
		log.Fatalf("Error initialising cookie instance: %v\n", err)
	}

	swiftContainer := os.Getenv("SWIFT_CONTAINER")
	secretKey := os.Getenv("AES_256_ENCRYPTION_KEY")
	swiftAPIUser := os.Getenv("SWIFT_API_USER")
	swiftAPIKey := os.Getenv("SWIFT_API_KEY")
	swiftAuthURL := os.Getenv("SWIFT_AUTH_URL")
	swiftTenantID := os.Getenv("SWIFT_TENANT_ID")
	siteAddress := os.Getenv("SITE_ADDRESS")
	// Initialise some things for start up
	swiftConn, err := filestore.InitSwift(
		secretKey,
		swiftAPIUser,
		swiftAPIKey,
		swiftAuthURL,
		swiftTenantID,
		siteAddress,
	)
	if err != nil {
		log.Fatalf("Error initialising swift sdk: %v\n", err)
	}

	appHMAC := apphmac.NewAppHMAC(secretKey)

	// Instantiate repositories
	andonRepository := repository.NewAndonRepository()
	andonIssueRepository := repository.NewAndonIssueRepository()
	authRepository := repository.NewAuthRepository()
	fileRepository := repository.NewFileRepository(swiftContainer, secretKey)
	galleryRepository := repository.NewGalleryRepository(secretKey, fileRepository)
	commentRepository := repository.NewCommentRepository(fileRepository)
	resourceRepository := repository.NewResourceRepository()
	serviceRepository := repository.NewServiceRepository()
	stockTrxRepository := repository.NewStockTransactionRepository()
	teamRepository := repository.NewTeamRepository()
	stockItemRepository := repository.NewStockItemRepository()
	userRepository := repository.NewUserRepository()
	searchRepository := repository.NewSearchRepository()

	// Instantiate services
	services := &router.Services{
		AndonService:            *service.NewAndonService(pgPool, swiftConn, andonRepository, commentRepository, galleryRepository),
		AndonIssueService:       *service.NewAndonIssueService(pgPool, andonIssueRepository),
		AuthService:             *service.NewAuthService(pgPool, authRepository),
		CommentService:          *service.NewCommentService(pgPool, swiftConn, commentRepository),
		FileService:             *service.NewFileService(pgPool, swiftConn, fileRepository),
		GalleryService:          *service.NewGalleryService(pgPool, swiftConn, appHMAC, fileRepository, galleryRepository),
		PDFService:              *service.NewPDFService(),
		ResourceService:         *service.NewResourceService(pgPool, commentRepository, galleryRepository, resourceRepository, serviceRepository),
		SearchService:           *service.NewSearchService(pgPool, userRepository, searchRepository),
		ServicesService:         *service.NewServicesService(pgPool, commentRepository, galleryRepository, resourceRepository, serviceRepository),
		StockItemService:        *service.NewStockItemService(pgPool, swiftConn, galleryRepository, stockItemRepository, commentRepository),
		StockTransactionService: *service.NewStockTransactionService(pgPool, stockTrxRepository),
		TeamService:             *service.NewTeamService(pgPool, teamRepository, userRepository),
		UserService:             *service.NewUserService(pgPool, userRepository),
	}

	// define server
	server := http.Server{
		Addr:    ":3000",
		Handler: router.NewRouter(services, *appHMAC),
	}

	// Initialising chromium instance for pdf generations
	pdf.InitChromium()
	defer pdf.ShutdownChromium()

	// Bind to a port and pass our router in
	fmt.Println("Local: 		https://localhost:3000")
	ip, err := localip.GetLocalIP()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("On Your Network:	https://%s:3000\n", ip)

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "dev" {
		err = server.ListenAndServeTLS("cert.pem", "key.pem")
	} else {
		err = server.ListenAndServe()
	}

	if err != nil {
		log.Fatalf("Error listening and serving: %v\n", err)
	}
}
