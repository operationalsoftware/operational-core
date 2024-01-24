package handlers

import (
	"net/http"
	"operationalcore/utils"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	var id int
	if cookie, err := r.Cookie("login-session"); err != nil {
		return
	} else {
		if err = utils.CookieInstance.Decode("login-session", cookie.Value, &id); err != nil {
			return
		}
	}

	// Delete cookie
	cookie := &http.Cookie{
		Name:     "login-session",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Unix(0, 0),
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, cookie)

	w.Header().Set("hx-redirect", "/login/password")
}
