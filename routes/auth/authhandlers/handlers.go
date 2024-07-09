package authhandlers

import (
	"app/db"
	"app/models/authmodel"
	"app/routes/auth/authviews"
	"app/utils"
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

func PasswordLogInPage(w http.ResponseWriter, r *http.Request) {
	_ = authviews.PasswordLoginPage(authviews.PasswordLoginPageProps{}).
		Render(w)

	return
}

func PasswordLogIn(w http.ResponseWriter, r *http.Request) {

	var err error
	retryPageProps := authviews.PasswordLoginPageProps{}

	_ = r.ParseForm()
	var formData authmodel.VerifyPasswordLoginInput
	err = utils.UnmarshalUrlValues(r.PostForm, &formData)
	if err != nil {
		retryPageProps.HasServerError = true
		retryPageProps.Username = formData.Username
		_ = authviews.PasswordLoginPage(retryPageProps).Render(w)
		return
	}

	ex := db.UseDB()
	var out authmodel.VerifyPasswordLoginOutput

	err = db.WithTx(ex, func(tx *sql.Tx) error {
		out, err = authmodel.VerifyPasswordLogin(tx, formData)
		return err
	})

	if err != nil {
		retryPageProps.HasServerError = true
		retryPageProps.Username = formData.Username
		_ = authviews.PasswordLoginPage(retryPageProps).Render(w)
		return
	}

	if out.FailureReason != "" {
		retryPageProps.LogInFailedError = out.FailureReason
		retryPageProps.Username = formData.Username
		_ = authviews.PasswordLoginPage(retryPageProps).Render(w)
		return
	}

	err = setSessionCookie(w, out.AuthUser.UserId)
	if err != nil {
		retryPageProps.HasServerError = true
		_ = authviews.PasswordLoginPage(retryPageProps).Render(w)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func Logout(w http.ResponseWriter, r *http.Request) {
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func setSessionCookie(w http.ResponseWriter, userID int) error {
	// set session cookie!
	encoded, err := utils.CookieInstance.Encode("login-session", fmt.Sprintf("%d", userID))
	if err != nil {
		return err
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

	return nil
}
