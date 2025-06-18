package pdfview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type PDFPageProps struct {
	Ctx       reqcontext.ReqContext
	Templates []string
}

func PDFGeneratorPage(p PDFPageProps) g.Node {

	content := g.Group([]g.Node{
		components.Card(
			g.El("form",
				g.Attr("method", "POST"),
				g.Attr("action", ""),
				g.Attr("target", "_blank"),
				h.Class("pdf-form"),
				g.Group([]g.Node{
					h.H1(g.Text("Generate PDF")),

					g.El("label", g.Text("PDF to test")),
					g.El("select",
						g.Attr("name", "template"),
						g.Group(g.Map(p.Templates, func(tmpl string) g.Node {
							return g.El("option",
								g.Attr("value", tmpl),
								g.Text(tmpl),
							)
						})),
					),

					g.El("label", g.Text("PDF Params (JSON)")),
					g.El("textarea",
						g.Attr("name", "params"),
					),

					components.Button(
						&components.ButtonProps{
							ButtonType: "Primary",
						},
						g.Text("Generate PDF"),
					),
				}),
			),
		),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Content: content,
		Title:   "PDF Templates",
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/pdfview/pdf_page.css"),
		},
	})

}
