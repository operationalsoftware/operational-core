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
		h.Button(
			h.Class("button"),
			h.Type("submit"),
			g.Text("Log In"),
		),
	})

	nfcButton := h.Button(
		h.Class("button nfc-login-button"),
		h.Type("button"),
		components.Icon(&components.IconProps{
			Identifier: "nfc-variant",
		}),
		g.Text("Log In with NFC"),
	)

	loginForm := func(content g.Node) g.Node {
		return components.Form(
			h.Method("POST"),
			h.ID("login-form"),
			h.Input(
				h.Type("hidden"),
				h.Name("EncryptedCredentials"),
				h.Value(encryptedCredentials),
			),
			content,
		)
	}

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

	showAllLoginOptionBtn := g.Group([]g.Node{
		components.Divider(
			h.Style("margin-top: 1.5rem; margin-bottom: 1.5rem;"),
			g.Text("OR"),
		),
		h.A(
			h.Href("/auth/password?ShowAll=true"),
			h.Button(
				h.Class("button show-all-login-button"),
				h.Type("button"),
				g.Text("Show all log in options"),
			),
		),
	})

	allLoginOptions := func() g.Node {
		return g.Group([]g.Node{
			loginForm(usernameAndPassword),
			components.Divider(
				h.Style("margin-top: 1.5rem; margin-bottom: 1.5rem;"),
				g.Text("OR"),
			),
			microsoft,
			components.Divider(
				h.Style("margin-top: 1.5rem; margin-bottom: 1.5rem;"),
				g.Text("OR"),
			),
			nfcButton,
			components.Divider(
				h.Style("margin-top: 1.5rem; margin-bottom: 1.5rem;"),
				g.Text("OR"),
			),
			qrcode,
		})
	}

	singleLoginOption := func(lastLoginMethod string) g.Node {
		var loginNode g.Node
		switch lastLoginMethod {
		case cookie.LoginMethodMicrosoft:
			loginNode = microsoft
		case cookie.LoginMethodNFC:
			loginNode = loginForm(nfcButton)
		case cookie.LoginMethodQRCODE:
			loginNode = qrcode
		default:
			loginNode = loginForm(usernameAndPassword)
		}

		return g.Group([]g.Node{
			loginNode,
			showAllLoginOptionBtn,
		})
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
				allLoginOptions(),
			),

			g.If(
				p.LastLoginMethod != "",
				singleLoginOption(p.LastLoginMethod),
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
		components.InlineScript("/static/js/app_nfc.js"),
		components.InlineScript("/internal/views/authview/password_login_page.js"),

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
