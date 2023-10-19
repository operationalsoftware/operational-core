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
	Text       string
	ButtonType ButtonType
	Size       ButtonSize
	Loading    bool
	Disabled   bool
}

func Button(p *ButtonProps) g.Node {
	classes := c.Classes{
		"primary": p.ButtonType == "", // default to primary
	}
	if p.ButtonType != "" {
		classes[string(p.ButtonType)] = true
	}
	if p.Size != "" {
		classes[string(p.Size)] = true
	}
	return h.Button(
		InlineStyle(Assets, "/Button.css"),
		g.If(p.Disabled || p.Loading, h.Disabled()),
		g.If(p.Loading, LoadingSpinner(LoadingSpinnerXs)),
		g.If(p.Loading, h.DataAttr("loading", "true")),
		classes,
		g.Text(p.Text),
	)
}
