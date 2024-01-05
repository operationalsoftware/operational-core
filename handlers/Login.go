package handlers

import (
	"net/http"
	"operationalcore/model"
	"operationalcore/partials"
	"operationalcore/utils"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Verify user
	user, err := model.VerifyUser(username, password)

	if err != nil {
		_ = partials.LoginForm(&partials.LoginFormProps{
			Error: err.Error(),
		}).Render(w)
	}

	if user.UserId != 0 {
		// Set cookie
		if encoded, err := utils.CookieInstance.Encode("operationalcore-session", user.UserId); err == nil {
			cookie := &http.Cookie{
				Name:     "operationalcore-session",
				Value:    encoded,
				HttpOnly: true,
				Secure:   true,
				Expires:  time.Now().Add(time.Hour * 24 * 30),
				Path:     "/",
				SameSite: http.SameSiteLaxMode,
			}

			http.SetCookie(w, cookie)
			w.Header().Set("hx-redirect", "/")
		} else {
			_ = partials.LoginForm(&partials.LoginFormProps{
				Error: "Something went wrong. Please try again later",
			}).Render(w)
		}
	}
}
