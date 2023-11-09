package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type ButtonSize string
type ButtonType string

const (
	ButtonSm ButtonSize = "small"
	ButtonLg ButtonSize = "large"
)

const (
	ButtonPrimary   ButtonType = "primary"
	ButtonDanger    ButtonType = "danger"
	ButtonSecondary ButtonType = "secondary"
	ButtonSuccess   ButtonType = "success"
	ButtonWarning   ButtonType = "warning"
)

type ButtonProps struct {
	Classes    c.Classes
	ButtonType ButtonType
	Size       ButtonSize
	Loading    bool
	Disabled   bool
}

func Button(p *ButtonProps, children ...g.Node) g.Node {
	if p.Classes == nil {
		p.Classes = c.Classes{}
	}
	if p.ButtonType != "" {
		p.Classes[string(p.ButtonType)] = true
	} else {
		p.Classes["primary"] = true
	}

	if p.Size != "" {
		p.Classes[string(p.Size)] = true
	}

	return h.Button(
		p.Classes,
		g.If(p.Disabled || p.Loading, h.Disabled()),
		g.If(p.Loading, LoadingSpinner(LoadingSpinnerSm)),
		g.If(p.Loading, h.DataAttr("loading", "true")),
		g.Group(children),
		InlineStyle(Assets, "/Button.css"),
	)
}
