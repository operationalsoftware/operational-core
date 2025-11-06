package authview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/cookie"
	"app/pkg/encryptcredentials"
	"app/pkg/reqcontext"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type PasswordLoginPageProps struct {
	Ctx              reqcontext.ReqContext
	Username         string
	LogInFailedError string
	HasServerError   bool
	LastLoginMethod  string
}

func PasswordLoginPage(p PasswordLoginPageProps) g.Node {
	encryptedCredentials := p.Ctx.Req.URL.Query().Get("EncryptedCredentials")

	decoded, _ := encryptcredentials.Decrypt(encryptedCredentials)

	usernameAndPassword := g.Group([]g.Node{
		components.Input(&components.InputProps{
			Label:       "Username",
			Name:        "Username",
			Placeholder: "Enter username",
			InputProps: []g.Node{
				h.Value(decoded.Username),
			},
		}),

		components.Input(&components.InputProps{
			Label:       "Password",
			Name:        "Password",
			InputType:   "password",
			Placeholder: "Enter password",
			InputProps: []g.Node{
				h.Value(decoded.Password),
			},
		}),

		h.Input(
			h.Type("hidden"),
			h.Name("EncryptedCredentials"),
			h.Value(encryptedCredentials),
		),

		h.Button(
			h.Class("button"),
			h.Type("submit"),
			g.Text("Log In"),
		),
	})

	nfc := h.Button(
		h.Class("button nfc-login-button"),
		h.Type("button"),
		components.Icon(&components.IconProps{
			Identifier: "nfc-variant",
		}),
		g.Text("Log In with NFC"),
	)

	microsoft := h.A(
		h.Class("button microsoft-login-link"),
		h.Href("/auth/microsoft/login"),
		components.Icon(&components.IconProps{
			Identifier: "microsoft-logo",
		}),
		g.Text("Log In with Microsoft"),
	)

	qrcode := h.A(
		h.Href("/auth/password/qrcode"),
		h.Button(
			h.Class("button qr-button"),
			h.Type("button"),
			components.Icon(&components.IconProps{
				Identifier: "qrcode",
			}),
			g.Text("Log In with QR Code"),
		),
	)

	baseForm := g.Group([]g.Node{
		components.Form(
			h.Method("POST"),
			h.ID("login-form"),
			usernameAndPassword,
		),

		components.Divider(
			h.Style("margin-top: 1.5rem; margin-bottom: 1.5rem;"),
			g.Text("OR"),
		),
		microsoft,
		components.Divider(
			h.Style("margin-top: 1.5rem; margin-bottom: 1.5rem;"),
			g.Text("OR"),
		),
		nfc,
		components.Divider(
			h.Style("margin-top: 1.5rem; margin-bottom: 1.5rem;"),
			g.Text("OR"),
		),
		qrcode,
	})

	lastLoginForm := []g.Node{}
	switch p.LastLoginMethod {
	case cookie.LOGIN_METHOD_PASSWORD:
		lastLoginForm = append(lastLoginForm, g.Group([]g.Node{
			components.Form(
				h.Method("POST"),
				h.ID("login-form"),
				usernameAndPassword,
			),
		}))
	case cookie.LOGIN_METHOD_MICROSOFT:
		lastLoginForm = append(lastLoginForm, microsoft)
	case cookie.LOGIN_METHOD_NFC:
		lastLoginForm = append(lastLoginForm, nfc)
	case cookie.LOGIN_METHOD_QRCODE:
		lastLoginForm = append(lastLoginForm, qrcode)
	}

	if p.LastLoginMethod != "" {
		lastLoginForm = append(lastLoginForm,
			components.Divider(
				h.Style("margin-top: 1.5rem; margin-bottom: 1.5rem;"),
				g.Text("OR"),
			),
			h.A(
				h.Href("/auth/password?ShowAll=1"),
				h.Button(
					h.Class("button show-all-login-button"),
					h.Type("button"),
					g.Text("Show all log in options"),
				),
			),
		)
	}

	content := g.Group([]g.Node{
		components.Card(

			h.Div(
				h.Class("logo-wrapper"),
				components.LogoImgFull(nil),
			),

			h.H1(g.Text("Welcome")),
			h.P(
				h.Style("margin-bottom: 1.5rem;"),
				g.Text("Please login to begin"),
			),

			g.If(
				p.LastLoginMethod == "",
				baseForm,
			),

			g.If(
				p.LastLoginMethod != "",
				g.Group(lastLoginForm),
			),

			g.If(
				p.LogInFailedError != "",
				h.P(
					h.Class("error"),
					g.Text(p.LogInFailedError),
				),
			),

			g.If(
				p.HasServerError,
				h.P(
					h.Class("error"),
					g.Text("Oops, something went wrong. Please try again and contact support if this issue persists"),
				),
			),
		),

		components.InlineStyle("/internal/views/authview/password_login_page.css"),
		g.If(
			// Script will be injected only on base form or when last login is nfc
			p.LastLoginMethod == "" || p.LastLoginMethod == cookie.LOGIN_METHOD_NFC,
			g.Group([]g.Node{
				components.InlineScript("/static/js/app_nfc.js"),
				components.InlineScript("/internal/views/authview/password_login_page.js"),
			}),
		),

		g.If(
			decoded.Username != "" && decoded.Password != "",
			h.Script(g.Raw(`
				document.addEventListener('DOMContentLoaded', function () {
					document.getElementById('auto-login-form').submit();
				});
			`)),
		),

		g.If(
			decoded.Username == "" && decoded.Password == "",
			h.Script(g.Raw(`
				document.addEventListener('DOMContentLoaded', function () {
					const url = new URL(window.location);
					url.searchParams.delete('EncryptedCredentials');
					window.history.replaceState({}, document.title, url);
				});
			`)),
		),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Content: content,
		Title:   "Log In",
	})
}
