package views

import (
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
			LinkPart: "/login",
			Icon:     "",
		},
	}

	loginContent := g.Group([]g.Node{

		h.H1(g.Text("Login Page")),

		partials.LoginForm(&partials.LoginFormProps{}),
	})

	return layout.Page(layout.PageProps{
		Title:   "Login",
		Content: loginContent,
		Crumbs:  crumbs,
		Ctx:     p.Ctx,
	})
}
