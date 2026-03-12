package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

// Deprecated: Use native HTML title text and standard positioning/CSS when needed.
type Position string

const (
	Top    Position = "top"
	Right  Position = "right"
	Bottom Position = "bottom"
	Left   Position = "left"
)

// Deprecated: Use the native HTML title attribute on the target element instead of this props wrapper.
type TooltipProps struct {
	Text     string
	Position Position
	Classes  c.Classes
}

// Deprecated: Use the native HTML title attribute on the target element instead of this tooltip wrapper component.
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
		InlineScript("/components/Tooltip.js"),
	)
}
