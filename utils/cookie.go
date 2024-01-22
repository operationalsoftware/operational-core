package utils

import (
	"os"

	"github.com/gorilla/securecookie"
)

var hashKey = []byte(os.Getenv("HASH_KEY"))
var blockKey = []byte(os.Getenv("BLOCK_KEY"))

var (
	CookieInstance = securecookie.New(
		hashKey,
		blockKey,
	)
)
