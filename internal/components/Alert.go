package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type AlertType string

const (
	AlertSuccess AlertType = "success"
	AlertError   AlertType = "error"
	AlertWarning AlertType = "warning"
)

type AlertProps struct {
	AlertType AlertType
	Message   string
}

func Alert(p *AlertProps) g.Node {
	classes := c.Classes{
		"alert": true,
	}

	if p.AlertType != "" {
		classes[string(p.AlertType)] = true
	} else {
		classes["success"] = true
	}

	return h.Div(
		classes,
		g.Text(p.Message),
		h.Button(
			h.Type("button"),
			h.Class("alert-close"),
			h.Data("alert-close", "true"),
			Icon(&IconProps{
				Identifier: "close",
			}),
		),
		InlineScript("/internal/components/Alert.js"),
	)
}
