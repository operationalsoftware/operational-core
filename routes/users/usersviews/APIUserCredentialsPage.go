package usersviews

import (
	"app/components"
	"app/layout"
	"app/reqcontext"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type APIUserCredentialsPageProps struct {
	Ctx      reqcontext.ReqContext
	Username string
	Password string
}

func APIUserCredentialsPage(p *APIUserCredentialsPageProps) g.Node {
	content := g.Group([]g.Node{

		apiUserCredentials(&apiUserCredentialsProps{
			username: p.Username,
			password: p.Password,
		}),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Success",
		Content: content,
	})
}

type apiUserCredentialsProps struct {
	username string
	password string
}

func apiUserCredentials(p *apiUserCredentialsProps) g.Node {
	return components.Card(
		h.Div(
			h.Class("api-user-credentials"),
			h.Div(
				h.Class("content"),
				h.H2(
					g.Text("API User Credentials"),
				),
				h.Span(g.Text("Username: ")),
				h.Span(g.Text(p.username)),
				h.Br(),
				h.Span(g.Text("Password: ")),
				h.Span(g.Text(p.password)),
			),
			components.Button(&components.ButtonProps{
				ButtonType: components.ButtonPrimary,
				Link:       "/users",
				Classes: c.Classes{
					"users-btn": true,
				},
			},
				g.Text("Back to Users"),
			),
		),
	)
}
