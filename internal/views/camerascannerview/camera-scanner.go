package camerascannerview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type CameraScannerAppProps struct {
	Ctx reqcontext.ReqContext
}

type AssetPaths struct {
	MainJs  string `json:"mainJs"`
	MainCss string `json:"mainCss"`
}

func CameraScannerApp(p *CameraScannerAppProps) g.Node {

	content := h.Div(
		h.ID("scanner-wrapper"),
		h.Class("fullscreen"),
		g.Attr("style", "display: none;"),
		h.Div(h.ID("scanner")),
		h.Button(
			h.ID("cancel-button"),
			h.Img(
				h.Src("/static/img/close.svg"), // Update the path as needed
				h.Alt("Close"),
				h.Width("24"),
				h.Height("24"),
			),
		),
		h.Div(
			h.ID("controls"),
			h.Select(
				h.ID("camera-select"),
				h.Option(
					h.Value(""),
					g.Text("Select Camera"),
				),
			),
		),
	)

	layoutMainPadding := false
	return layout.Page(layout.PageProps{
		Ctx:               p.Ctx,
		Title:             "QR Code Scanner",
		LayoutMainPadding: &layoutMainPadding,
		Content:           content,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/camerascannerview/camera-scanner.css"),
			h.Script(h.Src("/static/js/lib/html5-qrcode-2.3.8.min.js"), g.Attr("defer", "")),
			h.Script(h.Src("/static/js/camera-scanner.js"), g.Attr("defer", "")), // 🔥 defer added
		},
	})
}
