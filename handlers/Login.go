package handlers

import (
	"fmt"
	"net/http"
	"operationalcore/model"
	"operationalcore/partials"
	"operationalcore/utils"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("Username")
	password := r.FormValue("Password")

	// Verify user
	user, err := model.VerifyUser(username, password)

	if err != nil {
		_ = partials.LoginForm(&partials.LoginFormProps{
			Error: err.Error(),
		}).Render(w)
	}

	if user.UserId != 0 {
		// set cookie
		encoded, err := utils.CookieInstance.Encode("login-session", user.UserId)
		if err != nil {
			fmt.Println(err)
			_ = partials.LoginForm(&partials.LoginFormProps{
				Error: "An unexpected error occurred. Please report this issue.",
			}).Render(w)
			return
		}
		cookie := &http.Cookie{
			Name:     "login-session",
			Value:    encoded,
			HttpOnly: true,
			Secure:   true,
			Expires:  time.Now().Add(time.Hour * 24 * 30),
			Path:     "/",
			SameSite: http.SameSiteLaxMode,
		}

		http.SetCookie(w, cookie)
		w.Header().Set("hx-redirect", "/")
	}
}
