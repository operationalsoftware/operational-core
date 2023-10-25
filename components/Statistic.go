/*
<div class="stats-container">
<div class="stat-element">
<p class="stat-heading">Active users</p>
<div class="stat-value">
	<span>1,234</span>
</div>
</div>
<div class="stat-element">
<p class="stat-heading">Account Balance</p>
<div class="stat-value">
	<span>450000</span>
</div>
</div>
<div class="stat-element">
<p class="stat-heading">Withdraw Cost</p>
<div class="stat-value">
	{{ icon("account") }}
	<span>0.0001</span>
</div>
</div>
</div>
*/

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
		InlineStyle(Assets, "/Statistic.css"),
	)
}
