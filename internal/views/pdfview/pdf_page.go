package pdfview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/internal/pdftemplate"
	"app/pkg/printnode"
	"app/pkg/reqcontext"
	"fmt"
	"net/url"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type PDFPageProps struct {
	Ctx              reqcontext.ReqContext
	Templates        []pdftemplate.RegisteredTemplate
	SelectedTemplate *pdftemplate.RegisteredTemplate
}

func PDFGeneratorPage(p PDFPageProps) g.Node {
	exampleJSON := ""
	selectedTemplateName := ""
	if p.SelectedTemplate != nil {
		exampleJSON = p.SelectedTemplate.ExampleJSON
		selectedTemplateName = p.SelectedTemplate.Name
	}

	header := h.Section(
		h.Class("pdf-generator-header"),
		h.Nav(
			h.Class("pdf-generator-nav"),
			h.A(
				h.Href("/pdf/logs"),
				g.Text("PDF generation logs"),
			),
		),
		h.H1(g.Text("PDF Templates")),
	)

	templatesList := h.Section(
		h.Class("templates-section"),
		h.H2(g.Text("Document templates")),
		h.P(g.Text("View available templates, inspect example inputs, and load them into the PDF tester.")),
		h.Div(
			h.Class("templates-grid"),
			g.Group(g.Map(p.Templates, func(t pdftemplate.RegisteredTemplate) g.Node {
				return h.Div(
					h.Class("template-card"),
					h.Div(
						h.Class("template-card-header"),
						h.H3(g.Text(t.Name)),
						g.If(t.Description != "", h.P(h.Class("template-card-description"), g.Text(t.Description))),
					),
					h.Div(
						h.Class("template-card-body"),
						h.Div(
							h.Class("template-card-actions"),
							h.A(
								h.Class("button primary small"),
								h.Href(fmt.Sprintf("?TemplateName=%s", url.QueryEscape(t.Name))),
								g.Text("Load into tester"),
							),
						),
					),
				)
			})),
		),
	)

	form := h.Form(
		h.Method("POST"),
		h.Target("_blank"),
		g.Group([]g.Node{
			h.H2(g.Text("Test a PDF template")),
			h.Label(
				g.Text("PDF Template Name"),
				h.Select(
					h.Name("TemplateName"),
					g.Attr("onchange", "handleTemplateNameChange(event)"),
					h.Option(h.Value(""), h.Disabled(), g.If(selectedTemplateName == "", h.Selected()), g.Text("Select a template to begin")),
					g.Group(g.Map(p.Templates, func(t pdftemplate.RegisteredTemplate) g.Node {
						return h.Option(
							h.Value(t.Name),
							g.If(selectedTemplateName == t.Name, h.Selected()),
							g.Text(t.Name),
						)
					})),
				),
			),
			h.Label(
				g.Text("Template Input Data (JSON)"),
				h.Textarea(
					h.Name("InputData"),
					g.Text(exampleJSON),
				),
			),
			g.If(exampleJSON != "", h.Div(
				h.Class("example-input"),
				g.Text("Example Input"),
				h.Pre(g.Text(exampleJSON)),
			)),
			h.Button(
				h.Class("button primary"),
				h.Type("submit"),
				g.Text("Generate PDF"),
			),
		}),
	)

	return layout.Page(layout.PageProps{
		Ctx:   p.Ctx,
		Title: "PDF Templates",
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "text-box-outline",
				Title:          "PDFs",
				URLPart:        "pdf",
			},
			{
				Title:   "PDF Templates",
				URLPart: "generate",
			},
		},
		Content: g.Group([]g.Node{
			header,
			templatesList,
			form,
		}),
		AppendHead: []g.Node{components.InlineStyle("/internal/views/pdfview/pdf_page.css")},
		AppendBody: []g.Node{components.InlineScript("/internal/views/pdfview/pdf_page.js")},
	})
}

func printNodeStatusBox(status printnode.Status) g.Node {
	message := status.Message
	statusClass := "warning"

	if !status.Configured {
		message = "PrintNode API key is not configured. Set PRINTNODE_API_KEY to enable automated printing."
	} else if status.Reachable {
		statusClass = "success"
		switch {
		case status.AccountName != "" && status.AccountEmail != "":
			message = fmt.Sprintf("Connected to PrintNode as %s (%s).", status.AccountName, status.AccountEmail)
		case status.AccountEmail != "":
			message = fmt.Sprintf("Connected to PrintNode as %s.", status.AccountEmail)
		case status.AccountName != "":
			message = fmt.Sprintf("Connected to PrintNode as %s.", status.AccountName)
		default:
			message = "PrintNode connection is working."
		}
	} else {
		statusClass = "error"
		if message == "" {
			message = "Unable to reach PrintNode. Check the API key and network connectivity."
		}
	}

	return h.Div(
		c.Classes{
			"printnode-status-box": true,
			statusClass:            true,
		},
		h.Div(
			h.Class("printnode-status-title"),
			g.Text("PrintNode"),
		),
		h.P(g.Text(message)),
	)
}

func generationLogsSection(logs []model.PDFGenerationLog) g.Node {

	rows := components.TableRows{}

	for _, log := range logs {
		inputPreview := log.InputData
		if len(inputPreview) > 120 {
			inputPreview = inputPreview[:117] + "..."
		}
		documentCell := documentLinkCell(log.PDFTitle, log.FileURL)
		rows = append(rows, components.TableRow{
			Cells: []components.TableCell{
				{Contents: g.Text(log.TemplateName)},
				{Contents: documentCell},
				{Contents: h.Pre(g.Text(inputPreview))},
				{Contents: g.Text(log.CreatedAt.Format("2006-01-02 15:04:05"))},
			},
		})
	}

	return h.Section(
		h.Class("pdf-log-section"),
		h.H2(g.Text("Recent PDF generations")),
		components.Table(&components.TableProps{
			Columns: components.TableColumns{
				{TitleContents: g.Text("Template")},
				{TitleContents: g.Text("PDF title")},
				{TitleContents: g.Text("Input data")},
				{TitleContents: g.Text("Generated at")},
			},
			Rows: rows,
		}),
	)
}

func printLogsSection(
	logs []model.PDFPrintLog,
	printers []printnode.Printer,
	showAssignmentsNav bool,
	statusBox g.Node,
) g.Node {

	rows := components.TableRows{}

	for _, log := range logs {
		status := log.Status
		if status == "" {
			status = "pending"
		}
		printerLabel := log.PrinterName
		if printerLabel == "" {
			printerLabel = "-"
		}
		documentCell := documentLinkCell(log.PDFTitle, log.FileURL)
		rows = append(rows, components.TableRow{
			Cells: []components.TableCell{
				{Contents: g.Text(log.TemplateName)},
				{Contents: g.Text(log.RequirementName)},
				{Contents: g.Text(printerLabel)},
				{Contents: g.Text(fmt.Sprintf("%d", log.PrintNodeJobID))},
				{Contents: g.Text(status)},
				{Contents: g.Text(log.CreatedAt.Format("2006-01-02 15:04"))},
				{Contents: documentCell},
				{Contents: h.Form(
					h.Method("POST"),
					h.Action("/pdf/print"),
					h.Div(
						h.Class("print-log-action"),
						h.Input(
							h.Type("hidden"),
							h.Name("PrintLogID"),
							h.Value(fmt.Sprintf("%d", log.PDFPrintLogID)),
						),
						h.Select(
							h.Class("print-log-printer"),
							h.Name("PrinterName"),
							h.Option(h.Value(""), g.Text("Use logged printer")),
							g.Group(g.Map(printers, func(pr printnode.Printer) g.Node {
								return h.Option(h.Value(pr.Name), g.Text(pr.Name))
							})),
						),
						h.Button(
							h.Class("button secondary"),
							h.Type("submit"),
							g.Text("Reprint"),
						),
					),
				)},
			},
		})
	}

	return h.Section(
		h.Class("pdf-log-section"),
		g.If(showAssignmentsNav, h.Nav(
			h.Class("print-requirements-nav"),
			h.A(
				h.Href("/printing/printer-assignments"),
				g.Text("Printer assignments"),
			),
		)),
		h.Div(
			h.Class("print-jobs-titlebar"),
			h.H2(g.Text("Print jobs")),
			g.If(statusBox != nil, statusBox),
		),
		components.Table(&components.TableProps{
			Columns: components.TableColumns{
				{TitleContents: g.Text("Template")},
				{TitleContents: g.Text("Requirement")},
				{TitleContents: g.Text("Printer")},
				{TitleContents: g.Text("PrintNode Job")},
				{TitleContents: g.Text("Status")},
				{TitleContents: g.Text("Printed at")},
				{TitleContents: g.Text("Document")},
				{TitleContents: g.Text("Actions")},
			},
			Rows: rows,
		}),
	)
}

func documentLinkCell(title, fileURL string) g.Node {
	linkTitle := title
	if linkTitle == "" {
		linkTitle = "PDF Document"
	}
	if fileURL == "" {
		return g.Text("-")
	}
	return h.A(
		h.Href(fileURL),
		h.Target("_blank"),
		h.Rel("noopener noreferrer"),
		g.Text(linkTitle),
	)
}
