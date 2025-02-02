package cookie

import (
	"encoding/hex"
	"os"
	"sync"

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
