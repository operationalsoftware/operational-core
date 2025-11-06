package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"app/internal/model"
	"app/internal/service"
	"app/internal/views/authview"
	"app/pkg/appurl"
	"app/pkg/cookie"
	"app/pkg/encryptcredentials"
	"app/pkg/reqcontext"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

type AuthHandler struct {
	authService   service.AuthService
	msOauthConfig *oauth2.Config
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	h := &AuthHandler{authService: authService}

	// Initialize Microsoft OAuth config from environment if available
	clientID := os.Getenv("MS_OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("MS_OAUTH_SECRET")
	tenantID := os.Getenv("MS_OAUTH_TENANT_ID")
	redirectURL := "https://" + os.Getenv("SITE_ADDRESS") + "/auth/microsoft/callback"

	if clientID != "" && clientSecret != "" && tenantID != "" {
		h.msOauthConfig = &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "profile", "email", "User.Read"},
			Endpoint:     microsoft.AzureADEndpoint(tenantID),
		}
	}

	return h
}

type LoginPageURLVals struct {
	ShowAll bool
	Error   string
}

func (h *AuthHandler) PasswordLogInPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	props := authview.PasswordLoginPageProps{Ctx: ctx}
	lastLoginMethod := cookie.GetLastLoginMethod(r)

	var loginPageURLVals LoginPageURLVals
	err := appurl.Unmarshal(r.URL.Query(), &loginPageURLVals)
	if err != nil {
		props.HasServerError = true
		_ = authview.PasswordLoginPage(props).Render(w)
		return
	}

	if loginPageURLVals.ShowAll {
		cookie.ClearLastLoginCookie(w)
		lastLoginMethod = ""
	}

	props.LastLoginMethod = lastLoginMethod

	if loginPageURLVals.Error != "" {
		props.LogInFailedError = loginPageURLVals.Error
	}

	_ = authview.PasswordLoginPage(props).Render(w)
}

func (h *AuthHandler) PasswordLogIn(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	var err error
	retryPageProps := authview.PasswordLoginPageProps{Ctx: ctx}
	retryPageProps.LastLoginMethod = cookie.GetLastLoginMethod(r)

	err = r.ParseForm()
	if err != nil {
		retryPageProps.HasServerError = true
		_ = authview.PasswordLoginPage(retryPageProps).Render(w)
		return
	}

	var formData passwordLoginFormData
	err = appurl.Unmarshal(r.PostForm, &formData)

	attemptedMethod := cookie.LoginMethodPassword
	if formData.EncryptedCredentials != "" && formData.Username == "" && formData.Password == "" {
		attemptedMethod = cookie.LoginMethodNFC
	}

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

	cookie.SetLastLoginCookie(w, attemptedMethod)

	http.Redirect(w, r, "/", http.StatusSeeOther)
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
		duration = time.Duration(*out.AuthUser.SessionDurationMinutes) * time.Minute
	}

	err = cookie.SetSessionCookie(w, out.AuthUser.UserID, time.Duration(duration))
	if err != nil {
		retryPageProps.HasServerError = true
		_ = authview.QRcodeLoginPage(retryPageProps).Render(w)
		return
	}

	cookie.SetLastLoginCookie(w, cookie.LoginMethodQRCODE)

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

// Microsoft OAuth login: redirect to Microsoft
func (h *AuthHandler) MicrosoftLogin(w http.ResponseWriter, r *http.Request) {
	if h.msOauthConfig == nil {
		http.NotFound(w, r)
		return
	}

	// Generate a random state and set as short-lived cookie for CSRF protection
	var stateBytes [32]byte
	if _, err := rand.Read(stateBytes[:]); err != nil {
		http.Error(w, "failed to start login", http.StatusInternalServerError)
		return
	}
	state := base64.RawURLEncoding.EncodeToString(stateBytes[:])

	// Save state in a secure, short-lived cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth-state",
		Value:    state,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(1 * time.Hour),
	})

	// PKCE: generate code_verifier and send S256 code_challenge
	var verifierBytes [32]byte
	if _, err := rand.Read(verifierBytes[:]); err != nil {
		http.Error(w, "failed to start login", http.StatusInternalServerError)
		return
	}
	codeVerifier := base64.RawURLEncoding.EncodeToString(verifierBytes[:])
	sum := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.RawURLEncoding.EncodeToString(sum[:])

	// Save verifier in a secure cookie to use on callback
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth-pkce",
		Value:    codeVerifier,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(1 * time.Hour),
	})

	url := h.msOauthConfig.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	)
	http.Redirect(w, r, url, http.StatusFound)
}

// Microsoft OAuth callback: exchange code, fetch profile, find user by email, create session
func (h *AuthHandler) MicrosoftCallback(w http.ResponseWriter, r *http.Request) {
	if h.msOauthConfig == nil {
		http.NotFound(w, r)
		return
	}

	if errParam := r.URL.Query().Get("error"); errParam != "" {
		http.Error(w, "Microsoft sign-in failed: "+errParam, http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing authorization code", http.StatusBadRequest)
		return
	}

	// Verify state from cookie for CSRF protection
	stateQuery := r.URL.Query().Get("state")
	stateCookie, err := r.Cookie("oauth-state")
	if err != nil || stateCookie.Value == "" || stateCookie.Value != stateQuery {
		http.Error(w, "invalid OAuth state", http.StatusBadRequest)
		return
	}
	// Invalidate the state cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "oauth-state",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Unix(0, 0),
	})

	// Read PKCE verifier from cookie
	pkceCookie, _ := r.Cookie("oauth-pkce")
	var tokenOpts []oauth2.AuthCodeOption
	if pkceCookie != nil && pkceCookie.Value != "" {
		tokenOpts = append(tokenOpts, oauth2.SetAuthURLParam("code_verifier", pkceCookie.Value))
		// Invalidate pkce cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "oauth-pkce",
			Value:    "",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
			Expires:  time.Unix(0, 0),
		})
	}

	token, err := h.msOauthConfig.Exchange(r.Context(), code, tokenOpts...)
	if err != nil {
		log.Printf("Microsoft OAuth: token exchange failed: %v", err)
		http.Error(w, "failed to exchange token", http.StatusInternalServerError)
		return
	}

	client := h.msOauthConfig.Client(r.Context(), token)
	resp, err := client.Get("https://graph.microsoft.com/v1.0/me")
	if err != nil {
		http.Error(w, "failed to fetch user profile", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		http.Error(w, "failed to fetch user profile from Microsoft Graph", http.StatusBadGateway)
		return
	}

	var profile struct {
		Mail              string `json:"mail"`
		UserPrincipalName string `json:"userPrincipalName"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		http.Error(w, "failed to decode user profile", http.StatusInternalServerError)
		return
	}

	email := profile.Mail
	if email == "" {
		email = profile.UserPrincipalName
	}
	if email == "" {
		http.Error(w, "no email found on Microsoft account", http.StatusForbidden)
		return
	}

	authUser, err := h.authService.AuthenticateUserByEmail(r.Context(), email)
	if err != nil {
		log.Printf("Microsoft OAuth: AuthenticateUserByEmail failed: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	if authUser == nil {
		error := url.QueryEscape("Microsoft login allowed only for existing users")
		http.Redirect(w, r, fmt.Sprintf("/auth/password?Error=%s", error), http.StatusSeeOther)
		return
	}

	// Create session cookie using same logic as password login
	duration := cookie.DefaultSessionDurationMinutes
	if authUser.SessionDurationMinutes != nil {
		duration = time.Duration(*authUser.SessionDurationMinutes) * time.Minute
	}

	if err := cookie.SetSessionCookie(w, authUser.UserID, duration); err != nil {
		http.Error(w, "failed to set session", http.StatusInternalServerError)
		return
	}

	cookie.SetLastLoginCookie(w, cookie.LoginMethodMicrosoft)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
