package authviews

import (
	"app/components"
	"app/layout"
	"app/reqcontext"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type PasswordLoginPageProps struct {
	Ctx              reqcontext.ReqContext
	Username         string
	LogInFailedError string
	HasServerError   bool
}

func PasswordLoginPage(p PasswordLoginPageProps) g.Node {

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

				h.Label(
					g.Text("Username"),
					h.Input(
						h.Name("Username"),
						h.Value(p.Username),
						g.Attr("autocomplete", "on"),
					),
				),

				h.Label(
					g.Text("Password"),
					h.Input(
						h.Type("password"),
						h.Name("Password"),
					),
				),

				h.Button(
					h.Class("button"),
					h.Type("submit"),
					g.Text("Log In"),
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

		components.InlineStyle("/routes/auth/authviews/index.css"),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Content: content,
		Title:   "Log In",
	})
}
