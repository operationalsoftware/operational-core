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
	"app/pkg/tracker"
)

type AuthHandler struct {
	authService service.AuthService
	tracker     *tracker.Tracker
}

func NewAuthHandler(authService service.AuthService, tracker *tracker.Tracker) *AuthHandler {
	return &AuthHandler{authService: authService,
		tracker: tracker}
}

func (h *AuthHandler) PasswordLogInPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	_ = authview.PasswordLoginPage(authview.PasswordLoginPageProps{
		Ctx: ctx,
	}).
		Render(w)

}

func (h *AuthHandler) PasswordLogIn(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	var err error
	retryPageProps := authview.PasswordLoginPageProps{
		Ctx: ctx,
	}

	_ = r.ParseForm()
	var formData model.VerifyPasswordLoginInput
	err = appurl.Unmarshal(r.PostForm, &formData)
	if err != nil {
		retryPageProps.HasServerError = true
		retryPageProps.Username = formData.Username
		_ = authview.PasswordLoginPage(retryPageProps).Render(w)
		return
	}

	var out model.VerifyPasswordLoginOutput

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

	err = h.tracker.TrackEvent(r.Context(), tracker.TrackingEvent{
		UserID:     out.AuthUser.UserID,
		EventName:  "Auth.Login",
		OccurredAt: time.Now(),
		Context:    "PasswordLogIn",
		MetaData:   map[string]interface{}{"foo": "bar"},
	})
	if err != nil {
		fmt.Println(err)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) QRcodeLogInPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	_ = authview.QRcodeLoginPage(authview.QRcodeLoginPageProps{
		Ctx: ctx,
	}).
		Render(w)

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
}
