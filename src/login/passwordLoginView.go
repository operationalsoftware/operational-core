package login

import (
	"app/components"
	"app/layout"
	reqContext "app/reqcontext"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type passwordLoginViewProps struct {
	Ctx reqContext.ReqContext
}

func passwordLoginView(p *passwordLoginViewProps) g.Node {

	loginContent := g.Group([]g.Node{
		h.Div(
			h.Class("container"),
			h.H2(g.Text("Welcome")),
			h.P(g.Text("Please login to begin")),
			passwordLoginFormPartial(&passwordLoginFormPartialProps{}),
		),
	})

	return layout.Page(layout.PageProps{
		Title:   "Login",
		Content: loginContent,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/src/login/passwordLogin.css"),
		},
	})
}
