package views

import (
	"operationalcore/components"
	"operationalcore/layout"
	"operationalcore/partials"
	"operationalcore/utils"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type LoginProps struct {
	Ctx utils.Context
}

func Login(p *LoginProps) g.Node {
	var crumbs = []layout.Crumb{
		{
			Title:    "Login",
			LinkPart: "login",
			Icon:     "",
		},
	}

	loginContent := g.Group([]g.Node{
		h.Div(
			h.Class("container"),
			h.H2(g.Text("Welcome")),
			h.P(g.Text("Please login to begin")),
			partials.LoginForm(&partials.LoginFormProps{}),
		),
		components.InlineStyle(Assets, "/Login.css"),
	})

	return layout.Page(layout.PageProps{
		Title:   "Login",
		Content: loginContent,
		Crumbs:  crumbs,
		Ctx:     p.Ctx,
	})
}
