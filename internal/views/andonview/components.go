package andonview

import (
	"app/internal/components"
	"app/internal/model"
	"strconv"

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
	case model.AndonStatusRequiresAcknowledgement:
		classes["requires-acknowledgement"] = true
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

type acknowledgeButtonProps struct {
	andonID    int
	showText   bool
	buttonSize components.ButtonSize
	ReturnTo   string
}

func acknowledgeButton(p *acknowledgeButtonProps) g.Node {
	return components.Button(&components.ButtonProps{
		Classes: c.Classes{"acknowledge": true},
		Size:    p.buttonSize,
	},
		g.Attr("onclick", "updateAndon(event)"),
		g.Attr("data-id", strconv.Itoa(p.andonID)),
		g.Attr("data-action", "acknowledge"),
		g.Attr("title", "Acknowledge"),
		g.Attr("data-return-to", p.ReturnTo),

		components.Icon(&components.IconProps{
			Identifier: "gesture-tap-hold",
		}),

		g.If(p.showText, g.Text("Acknowledge")),
	)
}

type resolveButtonProps struct {
	andonID    int
	showText   bool
	buttonSize components.ButtonSize
	ReturnTo   string
}

func resolveButton(p *resolveButtonProps) g.Node {
	return components.Button(&components.ButtonProps{
		Classes: c.Classes{"resolve": true},
		Size:    p.buttonSize,
	},
		g.Attr("onclick", "updateAndon(event)"),
		g.Attr("data-id", strconv.Itoa(p.andonID)),
		g.Attr("data-action", "resolve"),
		g.Attr("title", "Resolve"),
		g.Attr("data-return-to", p.ReturnTo),

		components.Icon(&components.IconProps{
			Identifier: "check",
		}),

		g.If(p.showText, g.Text("Resolve")),
	)
}

type cancelButtonProps struct {
	andonID    int
	showText   bool
	buttonSize components.ButtonSize
	ReturnTo   string
}

func cancelButton(p *cancelButtonProps) g.Node {
	return components.Button(&components.ButtonProps{
		Classes: c.Classes{"danger": true},
		Size:    p.buttonSize,
	},
		g.Attr("onclick", "updateAndon(event)"),
		g.Attr("data-id", strconv.Itoa(p.andonID)),
		g.Attr("data-action", "cancel"),
		g.Attr("title", "Cancel"),
		g.Attr("data-return-to", p.ReturnTo),

		components.Icon(&components.IconProps{
			Identifier: "cancel",
		}),

		g.If(p.showText, g.Text("Cancel")),
	)
}

type reopenButtonProps struct {
	andonID    int
	showText   bool
	buttonSize components.ButtonSize
}

func reopenButton(p *reopenButtonProps) g.Node {
	return components.Button(&components.ButtonProps{
		Size: p.buttonSize,
	},
		g.Attr("onclick", "updateAndon(event)"),
		g.Attr("data-id", strconv.Itoa(p.andonID)),
		g.Attr("data-action", "reopen"),
		g.Attr("title", "Reopen"),

		components.Icon(&components.IconProps{
			Identifier: "restore",
		}),

		g.If(p.showText, g.Text("Reopen")),
	)
}
