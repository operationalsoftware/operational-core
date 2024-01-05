package views

import (
	"operationalcore/layout"
	"operationalcore/partials"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

var loginCrumb layout.Crumb = layout.Crumb{
	Text:     "Login",
	UrlToken: "login",
}

func Login() g.Node {
	crumbs := []layout.Crumb{
		loginCrumb,
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
