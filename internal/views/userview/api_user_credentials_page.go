package userview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
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
	return g.Group([]g.Node{
		components.Card(
			h.Div(
				h.H2(
					g.Text("API User Credentials"),
				),
				h.P(g.Text("Please keep these safe as you will not be able to access them again.")),
				h.Span(g.Text("Username: ")),
				h.Span(g.Text(p.username)),

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
		components.InlineStyle("/internal/views/userview/api_user_credentials_page.css"),
	})
}
