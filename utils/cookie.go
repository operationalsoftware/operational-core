package utils

import (
	"github.com/gorilla/securecookie"
)

var hashKey = []byte("1234567890123456789012345678901234567890123456789012345678901234")
var blockKey = []byte("12345678901234567890123456789012")

var (
	CookieInstance = securecookie.New(
		hashKey,
		blockKey,
	)
)
