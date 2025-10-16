package andonview

import (
	"app/internal/components"
	"app/internal/model"
	"strconv"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
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
		classes["work-in-progress"] = true
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
	andonID  int
	showText bool
	returnTo string
	isSmall  bool
}

func acknowledgeButton(p *acknowledgeButtonProps) g.Node {
	classes := c.Classes{
		"button":      true,
		"acknowledge": true,
		"small":       p.isSmall,
	}

	return h.Button(
		classes,
		g.Attr("onclick", "updateAndon(event)"),
		h.Title("Acknowledge"),
		h.Data("id", strconv.Itoa(p.andonID)),
		h.Data("action", "acknowledge"),
		g.If(p.returnTo != "", h.Data("return-to", p.returnTo)),

		components.Icon(&components.IconProps{
			Identifier: "gesture-tap-hold",
		}),

		g.If(p.showText, g.Text("Acknowledge")),
	)
}

type resolveButtonProps struct {
	andonID  int
	showText bool
	returnTo string
	isSmall  bool
}

func resolveButton(p *resolveButtonProps) g.Node {
	classes := c.Classes{
		"button":  true,
		"resolve": true,
		"small":   p.isSmall,
	}

	return h.Button(
		classes,
		g.Attr("onclick", "updateAndon(event)"),
		h.Title("Resolve"),
		h.Data("id", strconv.Itoa(p.andonID)),
		h.Data("action", "resolve"),
		g.If(p.returnTo != "", h.Data("return-to", p.returnTo)),

		components.Icon(&components.IconProps{
			Identifier: "check",
		}),

		g.If(p.showText, g.Text("Resolve")),
	)
}

type cancelButtonProps struct {
	andonID  int
	showText bool
	returnTo string
	isSmall  bool
}

func cancelButton(p *cancelButtonProps) g.Node {
	classes := c.Classes{
		"button": true,
		"danger": true,
		"small":  p.isSmall,
	}

	return h.Button(
		classes,
		g.Attr("onclick", "updateAndon(event)"),
		h.Title("Cancel"),
		h.Data("id", strconv.Itoa(p.andonID)),
		h.Data("action", "cancel"),
		g.If(p.returnTo != "", h.Data("return-to", p.returnTo)),

		components.Icon(&components.IconProps{
			Identifier: "cancel",
		}),

		g.If(p.showText, g.Text("Cancel")),
	)
}

type reopenButtonProps struct {
	andonID  int
	showText bool
	isSmall  bool
}

func reopenButton(p *reopenButtonProps) g.Node {
	classes := c.Classes{
		"button": true,
		"small":  p.isSmall,
	}

	return h.Button(
		classes,
		g.Attr("onclick", "updateAndon(event)"),
		h.Title("Reopen"),
		h.Data("id", strconv.Itoa(p.andonID)),
		h.Data("action", "reopen"),

		components.Icon(&components.IconProps{
			Identifier: "restore",
		}),

		g.If(p.showText, g.Text("Reopen")),
	)
}
