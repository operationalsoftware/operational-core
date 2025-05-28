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
	"app/pkg/cookie"
	"app/pkg/db"
	"app/pkg/env"
	"app/pkg/localip"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Load environment (if not in production or staging)
	err := env.Load()
	if err != nil {
		log.Fatalf("Error loading environment: %v\n", err)
	}

	// Verify environment variables
	err = env.Verify()
	if err != nil {
		log.Fatalf("Error verifying environment: %v\n", err)
	}

	// Create database connection (SQLite - to be removed)
	err = db.ConnectDB()
	if err != nil {
		log.Fatalf("Error connecting to SQLite: %v\n", err)
	}
	defer db.UseDB().Close()

	migrate.RunMigrations() // uses log.Fatal

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

	// Instantiate repositories
	authRepository := repository.NewAuthRepository()
	userRepository := repository.NewUserRepository()
	stockTrxRepository := repository.NewStockTransactionRepository()

	// Instantiate services
	services := &router.Services{
		AuthService:             *service.NewAuthService(pgPool, authRepository),
		UserService:             *service.NewUserService(pgPool, userRepository),
		SearchService:           *service.NewSearchService(pgPool, userRepository),
		StockTransactionService: *service.NewStockTransactionService(pgPool, stockTrxRepository),
	}

	// define server
	server := http.Server{
		Addr:    ":3000",
		Handler: router.NewRouter(services),
	}

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
