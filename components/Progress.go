package components

import (
	"strconv"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type ProgressType string

const (
	ProgressTypeSuccess ProgressType = "success"
	ProgressTypeWarning ProgressType = "warning"
	ProgressTypeDanger  ProgressType = "danger"
)

type ProgressProps struct {
	Percentage int
	Type       ProgressType
}

func Progress(p *ProgressProps, children ...g.Node) g.Node {
	classes := c.Classes{
		"progress-container": true,
	}

	if p.Percentage == 0 {
		p.Percentage = 30
	}

	if p.Type == "" {
		p.Type = ProgressTypeSuccess
	}

	return h.Div(
		classes,
		h.DataAttr("percentage", strconv.Itoa(p.Percentage)),
		h.DataAttr("type", string(p.Type)),
		h.Div(
			h.Class("progress-bar"),
		),
		h.Div(
			h.Class("progress"),
		),
		h.Span(
			h.Class("progress-label"),
			g.Text("0%"),
		),
		g.Group(children),
		InlineStyle("/components/Progress.css"),
		InlineScript("/components/Progress.js"),
	)
}
