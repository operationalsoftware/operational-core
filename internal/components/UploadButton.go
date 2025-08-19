package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type UploadButtonProps struct {
	ID      string
	Classes c.Classes
}

func UploadButton(p *UploadButtonProps) g.Node {
	if p.Classes == nil {
		p.Classes = c.Classes{}
	}

	p.Classes["upload-button"] = true

	if p.ID == "" {
		p.ID = "upload-button"
	}

	return Button(&ButtonProps{
		Size:     ButtonSm,
		Loading:  false,
		Disabled: false,
		Classes:  p.Classes,
	}, g.Group([]g.Node{
		Icon(&IconProps{
			Identifier: "upload",
		}),
		h.Label(h.For(p.ID), g.Text("Choose a file")),
		h.Input(
			g.Attr("name", "file"),
			g.Attr("type", "file"),
			g.Attr("id", p.ID),
		),
		h.Div(
			h.Class("file-info"),
			g.Text("No file chosen"),
		),
		InlineScript("/components/UploadButton.js"),
	}))
}
