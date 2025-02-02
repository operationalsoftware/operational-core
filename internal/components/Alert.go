package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
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
		h.Div(
			Icon(&IconProps{
				Identifier: "close",
			}),
		),
	)
}
