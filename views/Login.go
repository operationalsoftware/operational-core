package views

import (
	"operationalcore/layout"
	"operationalcore/partials"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

func Login() g.Node {
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

	return layout.Page(layout.PageParams{
		Title:   "Login",
		Content: loginContent,
		Crumbs:  crumbs,
	})
}
