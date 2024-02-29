package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
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
		InputLabel(&InputLabelProps{
			For: p.ID,
		}, g.Text("Choose a file")),
		h.Input(
			g.Attr("type", "file"),
			g.Attr("id", p.ID),
		),
		h.Div(
			h.Class("file-info"),
			g.Text("No file chosen"),
		),
		InlineStyle("/components/UploadButton.css"),
		InlineScript("/components/UploadButton.js"),
	}))
}
