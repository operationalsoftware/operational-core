package serviceview

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func StatusLegend() g.Node {
	type legendItem struct {
		class string
		text  string
	}

	items := []legendItem{
		{class: "threshold-80", text: ">= 80%"},
		{class: "threshold-90", text: ">= 90%"},
		{class: "is-due", text: "Is Due"},
	}

	return h.Div(
		h.Class("status-legend"),
		g.Map(items, func(item legendItem) g.Node {
			return h.Div(
				h.Span(
					h.Class("status-dot "+item.class),
				),
				g.Text(item.text),
			)
		}),
	)
}
