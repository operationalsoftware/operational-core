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
}

func Tooltip(p *TooltipProps, children ...g.Node) g.Node {
	classes := c.Classes{
		"tooltip": true,
		"top":     p.Position == "", // default to top
	}

	if p.Position != "" {
		classes[string(p.Position)] = true
	}

	return h.Div(
		classes,
		h.DataAttr("content", p.Text),
		g.Group(children),
		InlineStyle("/components/Tooltip.css"),
		InlineScript("/components/Tooltip.js"),
	)
}
