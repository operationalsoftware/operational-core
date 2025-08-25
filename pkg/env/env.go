package env

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/securecookie"
	"github.com/joho/godotenv"
)

// Load loads environment variables from .env file when not in production or
// staging
func Load() error {
	// Check if GO_ENV is "staging" or "production"
	goEnv := os.Getenv("GO_ENV")
	if goEnv != "staging" && goEnv != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Error loading .env file")
			return err
		}
	}

	return nil
}

// Verify verifies that all required environment variables are set
// and, if applicable, suggests a random value if not set
func Verify() error {
	var fail bool = false
	if os.Getenv("SECURE_COOKIE_HASH_KEY") == "" {
		fmt.Println("SECURE_COOKIE_HASH_KEY environment variable not set")
		secureCookieHashKey := securecookie.GenerateRandomKey(32)
		fmt.Println("Use SECURE_COOKIE_HASH_KEY=\"" + hex.EncodeToString(secureCookieHashKey) + "\"")
		fail = true
	}

	if os.Getenv("SECURE_COOKIE_BLOCK_KEY") == "" {
		fmt.Println("SECURE_COOKIE_BLOCK_KEY environment variable not set")
		secureCookieBlockKey := securecookie.GenerateRandomKey(32)
		fmt.Println("Use SECURE_COOKIE_BLOCK_KEY=\"" + hex.EncodeToString(secureCookieBlockKey) + "\"")
		fail = true
	}

	if os.Getenv("PG_USER") == "" {
		fmt.Println("PG_USER environment variable not set")
		fail = true
	}

	if os.Getenv("PG_PASSWORD") == "" {
		fmt.Println("PG_PASSWORD environment variable not set")
		fail = true
	}

	if os.Getenv("PG_HOST") == "" {
		fmt.Println("PG_HOST environment variable not set")
		fail = true
	}

	if os.Getenv("PG_PORT") == "" {
		fmt.Println("PG_PORT environment variable not set")
		fail = true
	}

	if os.Getenv("PG_DATABASE") == "" {
		fmt.Println("PG_DATABASE environment variable not set")
		fail = true
	}

	if os.Getenv("APP_ENV") == "" {
		fmt.Println("APP_ENV environment variable not set")
		fail = true
	}

	if os.Getenv("SWIFT_API_USER") == "" {
		fmt.Println("SWIFT_API_USER environment variable not set")
		fail = true
	}

	if os.Getenv("SWIFT_TENANT_ID") == "" {
		fmt.Println("SWIFT_TENANT_ID environment variable not set")
		fail = true
	}

	if os.Getenv("SWIFT_API_KEY") == "" {
		fmt.Println("SWIFT_API_KEY environment variable not set")
		fail = true
	}

	if os.Getenv("SWIFT_AUTH_URL") == "" {
		fmt.Println("SWIFT_AUTH_URL environment variable not set")
		fail = true
	}

	if os.Getenv("ORBIT_CONTAINER") == "" {
		fmt.Println("ORBIT_CONTAINER environment variable not set")
		fail = true
	}

	if fail {
		return fmt.Errorf("missing required environment variables")
	}

	return nil
}
