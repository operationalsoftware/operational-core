package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

// Deprecated: Use string CSS class names directly with h.Span and related HTML markup.
type BadgeSize string

// Deprecated: Use string CSS class names directly with h.Span and related HTML markup.
type BadgeType string

const (
	BadgeSm BadgeSize = "small"
	BadgeLg BadgeSize = "large"
)

const (
	BadgePrimary   BadgeType = "primary"
	BadgeDanger    BadgeType = "danger"
	BadgeSecondary BadgeType = "secondary"
	BadgeSuccess   BadgeType = "success"
	BadgeWarning   BadgeType = "warning"
)

// Deprecated: Build badge markup directly with h.Span and CSS classes instead of this props wrapper.
type BadgeProps struct {
	Classes c.Classes
	Size    BadgeSize
	Type    BadgeType
}

// Deprecated: Use h.Span directly and apply the "badge" CSS class with existing type/size classes instead of this wrapper component.
func Badge(p *BadgeProps, children ...g.Node) g.Node {

	if p.Classes == nil {
		p.Classes = c.Classes{}
	}

	if p.Type != "" {
		p.Classes[string(p.Type)] = true
	}
	if p.Size != "" {
		p.Classes[string(p.Size)] = true
	}

	p.Classes["badge"] = true

	return h.Span(
		p.Classes,
		g.Group(children),
	)
}
