package cookie

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/securecookie"
)

var (
	CookieInstance *securecookie.SecureCookie
	once           sync.Once
	err            error
)

func InitCookieInstance() error {
	once.Do(func() {

		var hashKey []byte
		hashKey, err = hex.DecodeString(os.Getenv("SECURE_COOKIE_HASH_KEY"))
		if err != nil {
			return
		}

		var blockKey []byte
		blockKey, err = hex.DecodeString(os.Getenv("SECURE_COOKIE_BLOCK_KEY"))
		if err != nil {
			return
		}

		cookieInstance := securecookie.New(
			hashKey,
			blockKey,
		)

		// Assign the connection to the package-level variable
		CookieInstance = cookieInstance
	})

	if err != nil {
		return err
	}

	return nil
}

func SetSessionCookie(w http.ResponseWriter, userID int, duration time.Duration) error {
	// set session cookie!
	fmt.Println("Setting cookie with expiration time: ", time.Now().Add(duration))

	encoded, err := CookieInstance.Encode("login-session", userID)
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:     "login-session",
		Value:    encoded,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(duration),
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)

	return nil
}
