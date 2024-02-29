package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type StatisticProps struct {
	Heading string
	Value   string
	Icon    string
	Classes c.Classes
}

func Statistic(p *StatisticProps) g.Node {
	if p.Classes == nil {
		p.Classes = c.Classes{}
	}

	p.Classes["statistic"] = true
	return h.Div(
		p.Classes,
		h.P(
			h.Class("stat-heading"),
			g.Text(p.Heading),
		),
		h.Div(
			h.Class("stat-value"),
			h.Span(
				g.Text(p.Value),
			),
		),
		InlineStyle("/components/Statistic.css"),
	)
}
