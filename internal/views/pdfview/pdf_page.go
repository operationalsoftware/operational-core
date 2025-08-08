package pdfview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/pdftemplate"
	"app/pkg/reqcontext"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type PDFPageProps struct {
	Ctx              reqcontext.ReqContext
	Templates        []pdftemplate.RegisteredTemplate
	SelectedTemplate *pdftemplate.RegisteredTemplate
}

func PDFGeneratorPage(p PDFPageProps) g.Node {

	exampleJSON := ""
	if p.SelectedTemplate != nil {
		exampleJSON = p.SelectedTemplate.ExampleJSON
	}

	content := h.FormEl(
		g.Attr("method", "POST"),
		g.Attr("action", ""),
		g.Attr("target", "_blank"),
		g.Group([]g.Node{
			h.H1(g.Text("Test a PDF template")),

			h.Label(
				g.Text("PDF Template Name"),

				h.Select(
					h.Name("TemplateName"),
					g.Attr("onchange", "handleTemplateNameChange(event)"),
					h.Option(h.Value(""), h.Disabled(), h.Selected(), g.Text("Select a template to begin")),
					g.Group(g.Map(p.Templates, func(t pdftemplate.RegisteredTemplate) g.Node {
						return h.Option(
							h.Value(t.Name),
							g.Text(t.Name),
						)
					})),
				),
			),

			h.Label(
				g.Text("Template Input Data (JSON)"),

				h.Textarea(
					h.Name("InputData"),
				),
			),

			g.If(exampleJSON != "",
				h.Div(h.Class("example-input"),
					g.Text("Example Input"),

					h.Pre(g.Text(exampleJSON)),
				),
			),

			components.Button(
				&components.ButtonProps{
					ButtonType: "primary",
				},
				g.Text("Generate PDF"),
			),
		}),
	)

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Content: content,
		Title:   "PDF Templates",
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/pdfview/pdf_page.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/pdfview/pdf_page.js"),
		},
	})

}
