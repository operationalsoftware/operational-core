package authview

import (
	"app/internal/components"
	"app/internal/layout"
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
}

func PasswordLoginPage(p PasswordLoginPageProps) g.Node {
	encryptedCredentials := p.Ctx.Req.URL.Query().Get("EncryptedCredentials")

	decoded, _ := encryptcredentials.Decrypt(encryptedCredentials)

	content := g.Group([]g.Node{
		components.Card(

			h.Div(
				h.Class("logo-wrapper"),
				components.LogoImgFull(nil),
			),

			h.H1(g.Text("Welcome")),
			h.P(g.Text("Please login to begin")),

			components.Form(
				h.Method("POST"),
				h.ID("login-form"),

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
				),

				h.Button(
					h.Class("button"),
					h.Type("submit"),
					g.Text("Log In"),
				),

				h.Div(
					h.Class("or-divider"),
					g.El("hr"),
					g.El("span", g.Text("OR")),
					g.El("hr"),
				),

				h.A(
					h.Class("button microsoft-login-link"),
					h.Href("/auth/microsoft/login"),
					components.Icon(&components.IconProps{
						Identifier: "microsoft-logo",
					}),
					g.Text("Log In with Microsoft"),
				),

				h.Div(
					h.Class("or-divider"),
					g.El("hr"),
					g.El("span", g.Text("OR")),
					g.El("hr"),
				),

				h.Button(
					h.Class("button nfc-login-button"),
					h.Type("button"),
					components.Icon(&components.IconProps{
						Identifier: "nfc-variant",
					}),
					g.Text("Log In with NFC"),
				),

				h.Div(
					h.Class("or-divider"),
					g.El("hr"),
					g.El("span", g.Text("OR")),
					g.El("hr"),
				),

				h.A(
					h.Href("/auth/password/qrcode"),
					h.Button(
						h.Class("button qr-button"),
						h.Type("button"),
						components.Icon(&components.IconProps{
							Identifier: "qrcode",
						}),
						g.Text("Log In with QR Code"),
					),
				),
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
