package handlers

import (
	"net/http"
	"operationalcore/utils"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	var id int
	if cookie, err := r.Cookie("operationalcore-session"); err != nil {
		return
	} else {
		if err = utils.CookieInstance.Decode("operationalcore-session", cookie.Value, &id); err != nil {
			return
		}
	}

	// Delete cookie
	cookie := &http.Cookie{
		Name:     "operationalcore-session",
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
