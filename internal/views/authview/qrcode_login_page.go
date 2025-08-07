package authview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type QRcodeLoginPageProps struct {
	Ctx              reqcontext.ReqContext
	Username         string
	LogInFailedError string
	HasServerError   bool
}

func QRcodeLoginPage(p QRcodeLoginPageProps) g.Node {
	encryptedCredentials := p.Ctx.Req.URL.Query().Get("EncryptedCredentials")

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
				h.ID("qrcode-login-form"),

				h.Div(
					h.Class("input-container"),

					h.Input(
						h.Class("qrcode-form-input"),
						h.Name("qrcode-input"),
						h.Type("text"),
						h.Placeholder("Scan QR Code"),
						h.Value(encryptedCredentials),
					),

					h.A(
						h.Class("camera-button"),
						h.Href("/camera-scanner?field=EncryptedCredentials"),
						components.Icon(&components.IconProps{
							Identifier: "camera",
						}),
					),
				),

				h.Button(
					h.Class("cancel-button"),
					h.Type("button"),
					g.Text("Cancel"),
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

		components.InlineStyle("/internal/views/authview/qrcode_login_page.css"),
		components.InlineScript("/internal/views/authview/qrcode_login_page.js"),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Content: content,
		Title:   "Log In",
	})
}
