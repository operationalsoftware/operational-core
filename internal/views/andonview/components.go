package andonview

import (
	"app/internal/components"
	"app/internal/model"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
)

func severityBadge(severity model.AndonSeverity, size components.BadgeSize) g.Node {

	return components.Badge(&components.BadgeProps{
		Size: size,
		Classes: c.Classes{
			"info":                  severity == model.AndonSeverityInfo,
			"self-resolvable":       severity == model.AndonSeveritySelfResolvable,
			"requires-intervention": severity == model.AndonSeverityRequiresIntervention,
		},
	},
		g.Text(string(severity)),
	)
}

func statusBadge(status model.AndonStatus, size components.BadgeSize) g.Node {
	classes := c.Classes{}

	switch status {
	case model.AndonStatusCancelled:
		classes["cancelled"] = true
	case model.AndonStatusClosed:
		classes["resolved"] = true
	case model.AndonStatusWorkInProgress:
		classes["wip"] = true
	case model.AndonStatusOutstanding:
		classes["outstanding"] = true
	}

	return components.Badge(&components.BadgeProps{
		Size:    size,
		Classes: classes,
	},
		g.Text(string(status)),
	)
}
