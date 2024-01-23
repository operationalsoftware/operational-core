package partials

import (
	"operationalcore/components"

	g "github.com/maragudk/gomponents"
	ghtmx "github.com/maragudk/gomponents-htmx"
	h "github.com/maragudk/gomponents/html"
)

type LoginFormProps struct {
	Error string
}

func LoginForm(p *LoginFormProps) g.Node {
	return components.Form(
		ghtmx.Post(""),
		ghtmx.Swap("outerHTML"),
		components.Input(&components.InputProps{
			Label:       "Username",
			Name:        "username",
			Placeholder: "Enter username",
		}),
		components.Input(&components.InputProps{
			Label:       "Password",
			Name:        "password",
			Placeholder: "Enter password",
			InputType:   "password",
		}),

		g.If(p.Error != "", components.Alert(&components.AlertProps{
			AlertType: components.AlertError,
			Message:   p.Error,
		})),

		components.Button(&components.ButtonProps{}, h.Type("submit"), g.Text("Login")),
	)
}
