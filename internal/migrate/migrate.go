package migrate

import (
	"app/pkg/db"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations() {
	ctx := context.Background()

	pgEnv := db.LoadPostgresEnv()

	defaultConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres",
		pgEnv.User, pgEnv.Password, pgEnv.Host, pgEnv.Port)

	defaultConn, err := pgx.Connect(ctx, defaultConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to default database: %v", err)
	}
	defer defaultConn.Close(ctx)

	// Check if database databaseExists
	databaseExists, err := checkDatabaseExists(ctx, defaultConn, pgEnv.Database)
	if err != nil {
		log.Fatalf("Error checking database existence: %v", err)
	}

	// Create database if missing
	if !databaseExists {
		log.Printf("Database %s does not exist. Creating...", pgEnv.Database)
		if err := createDatabase(ctx, defaultConn, pgEnv.Database); err != nil {
			log.Fatalf("Error creating database: %v", err)
		}
	}

	targetConnStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		pgEnv.User, pgEnv.Password, pgEnv.Host, pgEnv.Port, pgEnv.Database)
	// use a pool for compatibility with user service
	pgPool, err := pgxpool.New(context.Background(), targetConnStr)
	if err != nil {
		log.Fatalf("Unable to create Postgres connection pool: %v\n", err)
	}
	defer pgPool.Close()

	tx, err := pgPool.Begin(ctx)
	if err != nil {
		log.Fatalf("Failed to create transaction: %v\n", err)
	}
	defer tx.Rollback(ctx)

	// Check if database is initialised
	initialised, err := checkInitialised(ctx, pgPool)
	if err != nil {
		log.Fatalf("Error checking database initialisation: %v", err)
	}

	// Initialise if needed
	if !initialised {
		log.Println("Database is not initialised. Initialising...")
		if err := initialise(ctx, tx); err != nil {
			log.Fatalf("Error initialising database: %v", err)
		}
	}

	// Check if migrations are required
	migrationsRequired, err := checkMigrationRequired(ctx, tx)
	if err != nil {
		log.Fatalf("Error checking for required migrations: %v", err)
	}

	// Apply migrations if needed
	if migrationsRequired {
		log.Println("Applying database migrations...")

		if err := migrate(ctx, tx); err != nil {
			log.Fatalf("Error applying migrations: %v", err)
		}

		fmt.Println("Migration complete!")
	}

	// Commit the transaction
	if err := tx.Commit(ctx); err != nil {
		log.Fatalf("error committing transaction: %v", err)
	}
}
