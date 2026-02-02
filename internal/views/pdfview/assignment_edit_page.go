package pdfview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/printnode"
	"app/pkg/reqcontext"
	"strings"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type PrinterAssignmentEditPageProps struct {
	Ctx            reqcontext.ReqContext
	Requirement    string
	Printers       []printnode.Printer
	SelectedName   string
	PrintNodeReady bool
}

func PrinterAssignmentEditPage(p PrinterAssignmentEditPageProps) g.Node {
	selected := p.SelectedName
	options := []components.SearchSelectOption{}
	for _, pr := range p.Printers {
		options = append(options, components.SearchSelectOption{
			Text:     pr.Name,
			Value:    pr.Name,
			Selected: strings.EqualFold(pr.Name, selected),
		})
	}

	form := h.Form(
		h.Class("form"),
		h.Method("POST"),
		h.Action("/printing/printer-assignments"),
		h.Div(
			h.Label(g.Text("Printer")),
			components.SearchSelect(&components.SearchSelectProps{
				Name:        "PrinterName",
				Placeholder: "Select printer",
				Mode:        "single",
				Options:     options,
				Selected:    selected,
			}),
		),
		h.Input(
			h.Type("hidden"),
			h.Name("RequirementName"),
			h.Value(p.Requirement),
		),
		components.Button(
			&components.ButtonProps{
				ButtonType: components.ButtonPrimary,
			},
			h.Type("submit"),
			g.Text("Save"),
		),
	)

	return layout.Page(layout.PageProps{
		Ctx:   p.Ctx,
		Title: "Edit Printer Assignment",
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "printer-settings",
				Title:          "Printing",
				URLPart:        "printing",
			},
			{
				Title:   "Printer Assignments",
				URLPart: "printer-assignments",
			},
			{Title: p.Requirement},
		},
		Content: h.Section(
			h.Class("print-requirements-section printer-assignment-edit"),
			g.If(!p.PrintNodeReady, h.P(g.Text("PrintNode is not configured."))),
			g.If(p.PrintNodeReady, form),
		),
		AppendHead: []g.Node{components.InlineStyle("/internal/views/pdfview/pdf_page.css")},
	})
}
