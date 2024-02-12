package components

import (
	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type StatisticProps struct {
	Heading string
	Value   string
	Icon    string
}

func Statistic(p *StatisticProps) g.Node {
	return h.Div(
		h.Class("stat-element"),
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
