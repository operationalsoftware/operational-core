package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type UploadButtonProps struct {
	Id string
}

func UploadButton(p *UploadButtonProps) g.Node {
	return Button(&ButtonProps{
		Size:     ButtonSm,
		Loading:  false,
		Disabled: false,
		Classes: c.Classes{
			"upload-button": true,
		},
	}, g.Group([]g.Node{
		Icon(&IconProps{
			Identifier: "upload",
		}),
		InputLabel(&InputLabelProps{
			For: p.Id,
		}, g.Text("Choose a file")),
		h.Input(
			g.Attr("type", "file"),
			g.Attr("id", p.Id),
		),
		h.Div(
			h.Class("file-info"),
			g.Text("No file chosen"),
		),
		InlineStyle(Assets, "/UploadButton.css"),
		InlineScript(Assets, "/UploadButton.js"),
	}))
}
