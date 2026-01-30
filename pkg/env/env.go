package env

import (
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/gorilla/securecookie"
)

// Verify verifies that all required environment variables are set
// and, if applicable, suggests a random value if not set
func Verify() error {
	// Map of required environment variables; keys represent the required vars.
	required := map[string]struct{}{
		"SECURE_COOKIE_HASH_KEY":  {},
		"SECURE_COOKIE_BLOCK_KEY": {},
		"PG_USER":                 {},
		"PG_PASSWORD":             {},
		"PG_HOST":                 {},
		"PG_PORT":                 {},
		"PG_DATABASE":             {},
		"APP_ENV":                 {},
		"SWIFT_API_USER":          {},
		"SWIFT_TENANT_ID":         {},
		"SWIFT_API_KEY":           {},
		"SWIFT_AUTH_URL":          {},
		"SWIFT_CONTAINER":         {},
		"VAPID_PUBLIC_KEY":        {},
	}

	// Iterate in a stable order for consistent output
	keys := make([]string, 0, len(required))
	for k := range required {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	missing := make([]string, 0)

	for _, key := range keys {
		if os.Getenv(key) != "" {
			continue
		}

		switch key {
		case "SECURE_COOKIE_HASH_KEY":
			fmt.Println("SECURE_COOKIE_HASH_KEY environment variable not set")
			secureCookieHashKey := securecookie.GenerateRandomKey(32)
			fmt.Println("Use SECURE_COOKIE_HASH_KEY=\"" + hex.EncodeToString(secureCookieHashKey) + "\"")
		case "SECURE_COOKIE_BLOCK_KEY":
			fmt.Println("SECURE_COOKIE_BLOCK_KEY environment variable not set")
			secureCookieBlockKey := securecookie.GenerateRandomKey(32)
			fmt.Println("Use SECURE_COOKIE_BLOCK_KEY=\"" + hex.EncodeToString(secureCookieBlockKey) + "\"")
		default:
			fmt.Println(key + " environment variable not set")
		}

		missing = append(missing, key)
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}
	return nil
}

func IsProduction() bool {
	return os.Getenv("APP_ENV") == "production"
}
