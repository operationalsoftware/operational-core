package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
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
	Link       string
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

	p.Classes["button"] = true

	content := g.Group(
		[]g.Node{
			p.Classes,
			g.If(p.Disabled, h.Disabled()),
			g.If(p.Loading, LoadingSpinner(&LoadingSpinnerProps{
				Size: "sm",
			})),
			g.If(p.Loading, h.Data("loading", "true")),
			g.Group(children),
		},
	)

	var el g.Node
	if p.Link != "" {
		el = h.A(content, h.Href(p.Link))
	} else {
		el = h.Button(content)
	}

	return el
}
