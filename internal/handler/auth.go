package handler

import (
	"fmt"
	"net/http"
	"time"

	"app/internal/model"
	"app/internal/service"
	"app/internal/views/authview"
	"app/pkg/appurl"
	"app/pkg/cookie"
	"app/pkg/encryptcredentials"
	"app/pkg/reqcontext"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
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

	err = r.ParseForm()
	if err != nil {
		retryPageProps.HasServerError = true
		_ = authview.PasswordLoginPage(retryPageProps).Render(w)
		return
	}

	var formData passwordLoginFormData
	err = appurl.Unmarshal(r.PostForm, &formData)
	if err != nil {
		retryPageProps.HasServerError = true
		_ = authview.PasswordLoginPage(retryPageProps).Render(w)
		return
	}

	verifyLoginInput := model.VerifyPasswordLoginInput{}
	verifyLoginInput.Username = formData.Username
	verifyLoginInput.Password = formData.Password

	// decrypt credentials if encrypted credentials submitted
	if formData.EncryptedCredentials != "" {
		decodedData, err := encryptcredentials.Decrypt(formData.EncryptedCredentials)
		if err != nil {
			retryPageProps.HasServerError = true
			retryPageProps.LogInFailedError = ""
			_ = authview.PasswordLoginPage(retryPageProps).Render(w)
			return
		}
		verifyLoginInput.Username = decodedData.Username
		verifyLoginInput.Password = decodedData.Password
	}

	var out model.VerifyPasswordLoginOutput
	retryPageProps.Username = formData.Username

	out, err = h.authService.VerifyPasswordLogin(r.Context(), verifyLoginInput)

	if err != nil {
		retryPageProps.HasServerError = true
		_ = authview.PasswordLoginPage(retryPageProps).Render(w)
		return
	}

	if out.FailureReason != "" {
		retryPageProps.LogInFailedError = out.FailureReason
		_ = authview.PasswordLoginPage(retryPageProps).Render(w)
		return
	}

	duration := cookie.DefaultSessionDurationMinutes

	if out.AuthUser.SessionDurationMinutes != nil {
		duration = time.Duration(*out.AuthUser.SessionDurationMinutes) * time.Minute
	}

	err = cookie.SetSessionCookie(w, out.AuthUser.UserID, duration)
	if err != nil {
		retryPageProps.HasServerError = true
		_ = authview.PasswordLoginPage(retryPageProps).Render(w)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

type passwordLoginFormData struct {
	Username             string
	Password             string
	EncryptedCredentials string
}

func (h *AuthHandler) QRcodeLogInPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	_ = authview.QRcodeLoginPage(authview.QRcodeLoginPageProps{
		Ctx: ctx,
	}).
		Render(w)

	return
}

func (h *AuthHandler) QRcodeLogIn(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	var err error
	retryPageProps := authview.QRcodeLoginPageProps{
		Ctx: ctx,
	}

	encryptedString := r.FormValue("qrcode-input")
	if encryptedString == "" {
		retryPageProps.HasServerError = true
		retryPageProps.LogInFailedError = ""
		_ = authview.QRcodeLoginPage(retryPageProps).Render(w)
		return
	}

	decodedData, err := encryptcredentials.Decrypt(encryptedString)
	if err != nil {
		retryPageProps.HasServerError = true
		retryPageProps.LogInFailedError = ""
		_ = authview.QRcodeLoginPage(retryPageProps).Render(w)
		return
	}

	if decodedData.Username == "" && decodedData.Password == "" {
		retryPageProps.HasServerError = true
		retryPageProps.LogInFailedError = ""
		_ = authview.QRcodeLoginPage(retryPageProps).Render(w)
		return
	}

	var out model.VerifyPasswordLoginOutput

	out, err = h.authService.VerifyPasswordLogin(r.Context(), decodedData)

	if err != nil {
		retryPageProps.HasServerError = true
		retryPageProps.Username = decodedData.Username
		_ = authview.QRcodeLoginPage(retryPageProps).Render(w)
		return
	}

	if out.FailureReason != "" {
		retryPageProps.LogInFailedError = out.FailureReason
		retryPageProps.Username = decodedData.Username
		_ = authview.QRcodeLoginPage(retryPageProps).Render(w)
		return
	}

	duration := cookie.DefaultSessionDurationMinutes
	if out.AuthUser.SessionDurationMinutes != nil {
		duration = time.Duration(*out.AuthUser.SessionDurationMinutes)
	}

	err = cookie.SetSessionCookie(w, out.AuthUser.UserID, time.Duration(duration))
	if err != nil {
		retryPageProps.HasServerError = true
		_ = authview.QRcodeLoginPage(retryPageProps).Render(w)
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
