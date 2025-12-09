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
	Ctx               reqcontext.ReqContext
	Templates         []pdftemplate.RegisteredTemplate
	SelectedTemplate  *pdftemplate.RegisteredTemplate
	PrintNodeStatus   printnode.Status
	Printers          []printnode.Printer
	PrintRequirements []PrintRequirement
	GenerationLogs    []model.PDFGenerationLog
	PrintLogs         []model.PDFPrintLog
}

type PrintRequirement struct {
	Name               string
	Description        string
	SelectedPrinter    string
	DefaultPrinterName string
}

func PDFGeneratorPage(p PDFPageProps) g.Node {

	exampleJSON := ""
	selectedTemplateName := ""
	if p.SelectedTemplate != nil {
		exampleJSON = p.SelectedTemplate.ExampleJSON
		selectedTemplateName = p.SelectedTemplate.Name
	}

	printNodeMessage := p.PrintNodeStatus.Message
	alertType := components.AlertWarning

	if !p.PrintNodeStatus.Configured {
		printNodeMessage = "PrintNode API key is not configured. Set PRINTNODE_API_KEY to enable automated printing."
	} else if p.PrintNodeStatus.Reachable {
		alertType = components.AlertSuccess
		switch {
		case p.PrintNodeStatus.AccountName != "" && p.PrintNodeStatus.AccountEmail != "":
			printNodeMessage = fmt.Sprintf("Connected to PrintNode as %s (%s).", p.PrintNodeStatus.AccountName, p.PrintNodeStatus.AccountEmail)
		case p.PrintNodeStatus.AccountEmail != "":
			printNodeMessage = fmt.Sprintf("Connected to PrintNode as %s.", p.PrintNodeStatus.AccountEmail)
		case p.PrintNodeStatus.AccountName != "":
			printNodeMessage = fmt.Sprintf("Connected to PrintNode as %s.", p.PrintNodeStatus.AccountName)
		default:
			printNodeMessage = "PrintNode connection is working."
		}
	} else {
		alertType = components.AlertError
		if printNodeMessage == "" {
			printNodeMessage = "Unable to reach PrintNode. Check the API key and network connectivity."
		}
	}

	printNodeStatus := h.Div(
		h.Class("printnode-status"),
		components.Alert(&components.AlertProps{
			AlertType: alertType,
			Message:   printNodeMessage,
		}),
	)

	printRequirements := h.Section(
		h.Class("print-requirements-section"),
		h.H2(g.Text("Print requirements")),
		g.If(len(p.Printers) == 0,
			h.P(g.Text("Configure PrintNode to select printers for each requirement.")),
		),
		g.If(len(p.Printers) > 0,
			g.Group([]g.Node{
				h.P(h.Class("hint"), g.Text("Printer selections are saved locally in this browser.")),
				h.Div(
					h.Class("print-requirements-list"),
					g.Group(g.Map(p.PrintRequirements, func(req PrintRequirement) g.Node {
						return h.Div(
							h.Class("print-requirement-card"),
							h.Div(
								h.Class("print-requirement-header"),
								h.H3(g.Text(req.Name)),
								g.If(req.Description != "",
									h.P(h.Class("print-requirement-desc"), g.Text(req.Description)),
								),
							),
							h.Label(
								g.Text("Printer"),
								h.Select(
									h.Class("print-requirement-select"),
									g.Attr("data-requirement-name", req.Name),
									g.Attr("data-default-printer", req.SelectedPrinter),
									h.Option(
										h.Value(""),
										g.If(req.SelectedPrinter == "", h.Selected()),
										g.Text("Select a printer"),
									),
									g.Group(g.Map(p.Printers, func(printer printnode.Printer) g.Node {
										return h.Option(
											h.Value(fmt.Sprintf("%d", printer.ID)),
											g.If(req.SelectedPrinter == fmt.Sprintf("%d", printer.ID), h.Selected()),
											g.Text(printer.Name),
										)
									})),
								),
							),
							components.Button(
								&components.ButtonProps{
									ButtonType: components.ButtonPrimary,
									Classes:    c.Classes{"print-requirement-action": true},
								},
								g.Attr("data-requirement-name", req.Name),
								g.Text("Print with this requirement"),
							),
						)
					})),
				),
			}),
		),
	)

	printersList := g.If(p.PrintNodeStatus.Configured,
		h.Section(
			h.Class("printers-section"),
			h.H2(g.Text("Available printers")),
			g.If(len(p.Printers) == 0,
				h.P(g.Text("No printers found in PrintNode.")),
			),
			g.If(len(p.Printers) > 0,
				h.Div(
					h.Class("printers-list"),
					g.Group(g.Map(p.Printers, func(printer printnode.Printer) g.Node {
						computer := printer.ComputerName
						if computer == "" && printer.ComputerID != 0 {
							computer = fmt.Sprintf("ID %d", printer.ComputerID)
						}

						state := printer.State
						if state == "" {
							state = "Unknown"
						}

						return h.Div(
							h.Class("printer-card"),
							h.Div(
								h.Class("printer-card-header"),
								h.H3(g.Text(printer.Name)),
								g.If(printer.Default,
									h.Span(h.Class("badge"), g.Text("Default")),
								),
							),
							g.If(printer.Description != "",
								h.P(h.Class("printer-desc"), g.Text(printer.Description)),
							),
							h.Div(
								h.Class("printer-meta"),
								h.Div(
									h.Class("meta-row"),
									h.Span(h.Class("label"), g.Text("Computer")),
									h.Span(g.Text(computer)),
								),
								h.Div(
									h.Class("meta-row"),
									h.Span(h.Class("label"), g.Text("State")),
									h.Span(g.Text(state)),
								),
							),
						)
					})),
				),
			),
		),
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
						g.If(t.Description != "",
							h.P(h.Class("template-card-description"), g.Text(t.Description)),
						),
					),
					h.Div(
						h.Class("template-card-body"),
						g.If(t.ExampleJSON != "",
							h.Div(
								h.Class("template-card-example"),
								h.Div(h.Class("label"), g.Text("Example Input")),
								h.Pre(g.Text(t.ExampleJSON)),
							),
						),
						h.Div(
							h.Class("template-card-actions"),
							components.Button(
								&components.ButtonProps{
									ButtonType: components.ButtonSecondary,
									Link:       fmt.Sprintf("?TemplateName=%s", url.QueryEscape(t.Name)),
								},
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
			h.H1(g.Text("Test a PDF template")),

			h.Label(
				g.Text("PDF Template Name"),

				h.Select(
					h.Name("TemplateName"),
					g.Attr("onchange", "handleTemplateNameChange(event)"),
					h.Option(
						h.Value(""),
						h.Disabled(),
						g.If(selectedTemplateName == "", h.Selected()),
						g.Text("Select a template to begin"),
					),
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

	logRows := components.TableRows{}
	for _, log := range p.GenerationLogs {
		inputPreview := log.InputData
		if len(inputPreview) > 120 {
			inputPreview = inputPreview[:117] + "..."
		}
		documentCell := g.Text("-")
		if log.FileURL != "" {
			documentCell = h.A(h.Href(log.FileURL), g.Text(log.Filename))
		}
		logRows = append(logRows, components.TableRow{
			Cells: []components.TableCell{
				{Contents: g.Text(log.TemplateName)},
				{Contents: h.Pre(g.Text(inputPreview))},
				{Contents: documentCell},
				{Contents: g.Text(log.CreatedAt.Format("2006-01-02 15:04"))},
			},
		})
	}

	printRows := components.TableRows{}
	for _, log := range p.PrintLogs {
		status := log.Status
		if status == "" {
			status = "pending"
		}
		printerLabel := log.PrinterName
		if printerLabel == "" && log.PrinterID != 0 {
			printerLabel = fmt.Sprintf("Printer %d", log.PrinterID)
		}
		documentCell := g.Text("-")
		if log.FileURL != "" {
			documentCell = h.A(h.Href(log.FileURL), g.Text(log.Filename))
		}
		printRows = append(printRows, components.TableRow{
			Cells: []components.TableCell{
				{Contents: g.Text(log.TemplateName)},
				{Contents: g.Text(log.RequirementName)},
				{Contents: g.Text(printerLabel)},
				{Contents: g.Text(fmt.Sprintf("%d", log.PrintNodeJobID))},
				{Contents: g.Text(status)},
				{Contents: g.Text(log.CreatedAt.Format("2006-01-02 15:04"))},
				{Contents: documentCell},
				{Contents: h.Div(
					h.Class("print-log-action"),
					h.Select(
						h.Class("print-log-printer"),
						g.Attr("data-print-log-id", fmt.Sprintf("%d", log.ID)),
						h.Option(h.Value(""), g.Text("Use logged printer")),
						g.Group(g.Map(p.Printers, func(pr printnode.Printer) g.Node {
							label := pr.Name
							return h.Option(h.Value(fmt.Sprintf("%d", pr.ID)), g.Text(label))
						})),
					),
					components.Button(&components.ButtonProps{
						ButtonType: components.ButtonSecondary,
						Classes:    c.Classes{"print-log-reprint": true},
					},
						g.Attr("data-print-log-id", fmt.Sprintf("%d", log.ID)),
						g.Text("Reprint"),
					),
				)},
			},
		})
	}

	return layout.Page(layout.PageProps{
		Ctx: p.Ctx,
		Content: g.Group([]g.Node{
			printNodeStatus,
			templatesList,
			form,
			printersList,
			printRequirements,
			h.Section(
				h.Class("pdf-log-section"),
				h.H2(g.Text("Recent PDF generations")),
				g.If(len(p.GenerationLogs) == 0,
					h.P(g.Text("No PDF generations logged yet.")),
				),
				g.If(len(p.GenerationLogs) > 0,
					components.Table(
						&components.TableProps{
							Columns: components.TableColumns{
								{TitleContents: g.Text("Template")},
								{TitleContents: g.Text("Input data")},
								{TitleContents: g.Text("Document")},
								{TitleContents: g.Text("Generated at")},
							},
							Rows: logRows,
						},
					),
				),
			),
			h.Section(
				h.Class("pdf-log-section"),
				h.H2(g.Text("Recent print jobs")),
				g.If(len(p.PrintLogs) == 0,
					h.P(g.Text("No print jobs logged yet.")),
				),
				g.If(len(p.PrintLogs) > 0,
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
						Rows: printRows,
					}),
				),
			),
		}),
		Title: "PDF Templates",
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/pdfview/pdf_page.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/pdfview/pdf_page.js"),
		},
	})

}
