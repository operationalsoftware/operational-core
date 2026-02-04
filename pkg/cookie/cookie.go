package cookie

import (
	"encoding/hex"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/securecookie"
)

var (
	CookieInstance *securecookie.SecureCookie
	once           sync.Once
	err            error
)

const (
	LoginMethodPassword  = "password"
	LoginMethodMicrosoft = "microsoft"
	LoginMethodNFC       = "nfc"
	LoginMethodQRCODE    = "qrcode"
)

const DefaultSessionDurationMinutes = time.Hour * 24 * 30
const pushSubscriptionCookieName = "push-subscription"
const pushRebindCookieName = "push-rebind"

type SessionData struct {
	UserID    int
	ExpiresAt time.Time
}

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
	cookieDate := SessionData{
		UserID:    userID,
		ExpiresAt: time.Now().Add(duration),
	}

	encoded, err := CookieInstance.Encode("login-session", cookieDate)
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

func GetSessionData(r *http.Request) (SessionData, error) {
	sessCookie, err := r.Cookie("login-session")
	if err != nil {
		return SessionData{}, err
	}

	var data SessionData
	if err := CookieInstance.Decode("login-session", sessCookie.Value, &data); err != nil {
		return SessionData{}, err
	}

	return data, nil
}

func SetLastLoginCookie(w http.ResponseWriter, method string) {
	cookie := &http.Cookie{
		Name:     "last-login-method",
		Value:    method,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)
}

func ClearLastLoginCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "last-login-method",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}

func SetPushSubscriptionCookie(w http.ResponseWriter, endpoint string) error {
	if strings.TrimSpace(endpoint) == "" {
		return nil
	}

	encoded, err := CookieInstance.Encode(pushSubscriptionCookieName, endpoint)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     pushSubscriptionCookieName,
		Value:    encoded,
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(DefaultSessionDurationMinutes),
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)

	return nil
}

func GetPushSubscriptionEndpoint(r *http.Request) (string, error) {
	endpointCookie, err := r.Cookie(pushSubscriptionCookieName)
	if err != nil {
		return "", err
	}

	var endpoint string
	if err := CookieInstance.Decode(pushSubscriptionCookieName, endpointCookie.Value, &endpoint); err != nil {
		return "", err
	}

	return endpoint, nil
}

func ClearPushSubscriptionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     pushSubscriptionCookieName,
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}

func SetPushRebindCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     pushRebindCookieName,
		Value:    "1",
		HttpOnly: false,
		Secure:   true,
		Expires:  time.Now().Add(10 * time.Minute),
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearPushRebindCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     pushRebindCookieName,
		Value:    "",
		HttpOnly: false,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}

func GetLastLoginMethod(r *http.Request) string {
	cookie, err := r.Cookie("last-login-method")
	if err != nil {
		return ""
	}

	switch cookie.Value {
	case LoginMethodPassword, LoginMethodMicrosoft, LoginMethodNFC, LoginMethodQRCODE:
		return cookie.Value
	default:
		return ""
	}
}
