package components

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type ToastProps struct {
	Contents g.Node
	Type     string // e.g., "success", "error", "info"
}

func Toast(p *ToastProps) g.Node {
	class := "toast"
	switch p.Type {
	case "success":
		class += " toast-success"
	case "error":
		class += " toast-error"
	default:
		class += " toast-info"
	}

	return h.Div(
		h.Class(class),

		p.Contents,

		InlineScript("/internal/components/Toast.js"),
	)
}
