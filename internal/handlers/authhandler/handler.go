package authhandler

import (
	"net/http"
	"time"

	"app/internal/models"
	"app/internal/services/authservice"
	"app/internal/views/authview"
	"app/pkg/cookie"
	"app/pkg/reqcontext"
	"app/pkg/urlvalues"
)

type AuthHandler struct {
	authService authservice.AuthService
}

func NewAuthHandler(authService authservice.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) PasswordLogInPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	_ = authview.PasswordLoginPage(authview.PasswordLoginPageProps{
		Ctx: ctx,
	}).
		Render(w)

	return
}

func (h *AuthHandler) PasswordLogIn(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	var err error
	retryPageProps := authview.PasswordLoginPageProps{
		Ctx: ctx,
	}

	_ = r.ParseForm()
	var formData models.VerifyPasswordLoginInput
	err = urlvalues.Unmarshal(r.PostForm, &formData)
	if err != nil {
		retryPageProps.HasServerError = true
		retryPageProps.Username = formData.Username
		_ = authview.PasswordLoginPage(retryPageProps).Render(w)
		return
	}

	var out models.VerifyPasswordLoginOutput

	out, err = h.authService.VerifyPasswordLogin(r.Context(), formData)

	if err != nil {
		retryPageProps.HasServerError = true
		retryPageProps.Username = formData.Username
		_ = authview.PasswordLoginPage(retryPageProps).Render(w)
		return
	}

	if out.FailureReason != "" {
		retryPageProps.LogInFailedError = out.FailureReason
		retryPageProps.Username = formData.Username
		_ = authview.PasswordLoginPage(retryPageProps).Render(w)
		return
	}

	err = setSessionCookie(w, out.AuthUser.UserID)
	if err != nil {
		retryPageProps.HasServerError = true
		_ = authview.PasswordLoginPage(retryPageProps).Render(w)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
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
	encoded, err := cookie.CookieInstance.Encode("login-session", userID)
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
