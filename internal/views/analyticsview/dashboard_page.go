package analyticsview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"
	"app/pkg/tracker"
	"fmt"
	"strconv"
	"time"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type DashboardPageProps struct {
	Ctx   reqcontext.ReqContext
	Stats tracker.EventStats
}

// func DashboardPage(p DashboardPageProps) g.Node {

// 	content := g.Group([]g.Node{
// 		h.H1(g.Text("Heading")),
// 	})

// 	return layout.Page(layout.PageProps{
// 		Ctx:     p.Ctx,
// 		Content: content,
// 		Title:   "Log In",
// 	})
// }

func DashboardPage(p DashboardPageProps) g.Node {
	content := h.Div(h.Class("container"),
		h.Header(h.Class("header"),
			h.H1(g.Text("📊 User Events Analytics")),
			h.P(h.Class("subtitle"), g.Text("Real-time insights into user behavior")),
		),

		// Summary Cards
		h.Div(h.Class("summary-grid"),
			summaryCard("Total Events", strconv.Itoa(p.Stats.TotalEvents), "🎯"),
			summaryCard("Unique Users", strconv.Itoa(p.Stats.UniqueUsers), "👥"),
			summaryCard("Avg Events/User", fmt.Sprintf("%.1f", float64(p.Stats.TotalEvents)/float64(max(p.Stats.UniqueUsers, 1))), "📈"),
			summaryCard("Active Today", "Loading...", "⚡"),
		),

		// Charts Section
		h.Div(h.Class("charts-grid"),
			// Top Events Chart
			h.Div(h.Class("chart-card"),
				h.H3(g.Text("Top Events")),
				h.Div(h.Class("bar-chart"), h.ID("top-events-chart"),
					topEventsChart(p.Stats.TopEvents),
				),
			),

			// Events Over Time
			h.Div(h.Class("chart-card"),
				h.H3(g.Text("Events Over Time (Last 7 Days)")),
				h.Div(h.Class("line-chart"), h.ID("timeline-chart"),
					timelineChart(p.Stats.EventsOverTime),
				),
			),

			// Hourly Distribution
			h.Div(h.Class("chart-card"),
				h.H3(g.Text("Events by Hour")),
				h.Div(h.Class("bar-chart horizontal"), h.ID("hourly-chart"),
					hourlyChart(p.Stats.EventsByHour),
				),
			),
		),

		// Recent Events Table
		h.Div(h.Class("table-card"),
			h.H3(g.Text("Recent Events")),
			h.Div(h.Class("table-container"),
				recentEventsTable(p.Stats.RecentEvents),
			),
		),

		// Footer
		h.Footer(h.Class("footer"),
			h.P(g.Text("Dashboard last updated: "), h.Span(h.ID("last-updated"), g.Text(time.Now().Format("15:04:05")))),
		),
	)

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Content: content,
		Title:   "Log In",
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/analyticsview/dashboard_page.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/analyticsview/dashboard_page.js"),
		},
	})

}

func summaryCard(title, value, icon string) g.Node {
	return h.Div(h.Class("summary-card"),
		h.Div(h.Class("card-icon"), g.Text(icon)),
		h.Div(h.Class("card-content"),
			h.H3(g.Text(value)),
			h.P(g.Text(title)),
		),
	)
}

func topEventsChart(events []tracker.EventCount) g.Node {
	if len(events) == 0 {
		return h.P(h.Class("no-data"), g.Text("No events data available"))
	}

	maxCount := 0
	for _, e := range events {
		if e.Count > maxCount {
			maxCount = e.Count
		}
	}

	var bars []g.Node
	for _, event := range events {
		percentage := float64(event.Count) / float64(maxCount) * 100
		bars = append(bars,
			h.Div(h.Class("bar-item"),
				h.Div(h.Class("bar-label"), g.Text(event.EventName)),
				h.Div(h.Class("bar-container"),
					h.Div(h.Class("bar"),
						// h.Style(fmt.Sprintf("width: %.1f%%", percentage)),
						g.Attr("style", fmt.Sprintf("width: %.1f%%", percentage)),
					),
					h.Span(h.Class("bar-value"), g.Text(strconv.Itoa(event.Count))),
				),
			),
		)
	}

	return h.Div(bars...)
}

func timelineChart(points []tracker.TimeSeriesPoint) g.Node {
	if len(points) == 0 {
		return h.P(h.Class("no-data"), g.Text("No timeline data available"))
	}

	maxCount := 0
	for _, p := range points {
		if p.Count > maxCount {
			maxCount = p.Count
		}
	}

	var chartPoints []g.Node
	for i, point := range points {
		height := float64(point.Count) / float64(maxCount) * 100
		left := float64(i) / float64(len(points)-1) * 100

		chartPoints = append(chartPoints,
			h.Class("timeline-container"),
			h.Div(h.Class("timeline-point"),
				// h.Style(fmt.Sprintf("left: %.1f%%; height: %.1f%%", left, height)),
				g.Attr("style", fmt.Sprintf("left: %.1f%%; height: %.1f%%", left, height)),
				g.Attr("title", fmt.Sprintf("%s: %d events", point.Date, point.Count)),
				// h.Title(fmt.Sprintf("%s: %d events", point.Date, point.Count)),
			),
		)
	}

	return h.Div(chartPoints...)
}

func hourlyChart(hours []tracker.HourlyCount) g.Node {
	if len(hours) == 0 {
		return h.P(h.Class("no-data"), g.Text("No hourly data available"))
	}

	maxCount := 0
	for _, h := range hours {
		if h.Count > maxCount {
			maxCount = h.Count
		}
	}

	var bars []g.Node
	for _, hour := range hours {
		percentage := float64(hour.Count) / float64(maxCount) * 100
		bars = append(bars,
			h.Div(h.Class("hourly-bar"),
				h.Div(h.Class("hourly-label"), g.Text(fmt.Sprintf("%02d:00", hour.Hour))),
				h.Div(h.Class("hourly-bar-fill"),
					// h.Style(fmt.Sprintf("width: %.1f%%", percentage)),
					g.Attr("style", fmt.Sprintf("width: %.1f%%", percentage)),
				),
				h.Span(h.Class("hourly-value"), g.Text(strconv.Itoa(hour.Count))),
			),
		)
	}

	return h.Div(bars...)
}

func recentEventsTable(events []tracker.RecentEvent) g.Node {
	if len(events) == 0 {
		return h.P(h.Class("no-data"), g.Text("No recent events"))
	}

	var rows []g.Node
	rows = append(rows,
		h.Class("events-table"),
		h.Tr(h.Class("table-header"),
			h.Th(g.Text("Event Name")),
			h.Th(g.Text("User ID")),
			h.Th(g.Text("Time")),
			h.Th(g.Text("Context")),
		),
	)

	for _, event := range events {
		rows = append(rows,
			h.Tr(
				h.Td(h.Class("event-name"), g.Text(event.EventName)),
				h.Td(g.Text(strconv.Itoa(event.UserID))),
				h.Td(h.Class("timestamp"), g.Text(event.OccurredAt.Format("15:04:05"))),
				h.Td(h.Class("context"), g.Text(truncateString(event.Context, 50))),
			),
		)
	}

	return h.Table(rows...)
}

func truncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
