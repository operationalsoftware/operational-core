package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type Position string

const (
	Top    Position = "top"
	Right  Position = "right"
	Bottom Position = "bottom"
	Left   Position = "left"
)

type TooltipProps struct {
	Text     string
	Position Position
	Classes  c.Classes
}

func Tooltip(p *TooltipProps, children ...g.Node) g.Node {
	if p.Classes == nil {
		p.Classes = c.Classes{}
	}

	if p.Position != "" {
		p.Classes[string(p.Position)] = true
	}

	p.Classes["tooltip"] = true
	p.Classes["top"] = p.Position == ""

	return h.Div(
		p.Classes,
		h.DataAttr("content", p.Text),
		g.Group(children),
		InlineStyle("/components/Tooltip.css"),
		InlineScript("/components/Tooltip.js"),
	)
}
