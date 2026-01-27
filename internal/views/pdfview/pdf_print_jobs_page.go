package pdfview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/printnode"
	"app/pkg/reqcontext"
	"fmt"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type PDFPrintJobsPageProps struct {
	Ctx             reqcontext.ReqContext
	PrintLogs       []model.PDFPrintLog
	Printers        []printnode.Printer
	PrintNodeStatus printnode.Status
}

func printLogsSection(
	logs []model.PDFPrintLog,
	printers []printnode.Printer,
	showAssignmentsNav bool,
	printNodeStatus printnode.Status,
) g.Node {
	rows := components.TableRows{}

	for _, log := range logs {
		requirementName := log.RequirementName

		printNodeJobId := "-"
		if log.PrintNodeJobID != nil {
			printNodeJobId = fmt.Sprintf("%d", *log.PrintNodeJobID)
		}
		errorMessage := "-"
		hasError := false
		if log.ErrorMessage != nil && *log.ErrorMessage != "" {
			errorMessage = *log.ErrorMessage
			hasError = true
		}
		documentCell := documentLinkCell(log.PDFTitle, log.FileURL)

		printerOverrideForm := h.Form(
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
					h.Option(h.Value(""), g.Text("Use assigned printer")),
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
		)

		rows = append(rows, components.TableRow{
			Classes: c.Classes{
				"print-log-error-row": hasError,
			},
			Cells: []components.TableCell{
				{Contents: g.Text(log.TemplateName)},
				{Contents: g.Text(requirementName)},
				{Contents: g.Text(printNodeJobId)},
				{Contents: documentCell},
				{Contents: g.Text(log.CreatedByUsername)},
				{Contents: g.Text(log.CreatedAt.Format("2006-01-02 15:04"))},
				{Contents: h.Pre(g.Text(errorMessage))},
				{Contents: printerOverrideForm},
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
			printNodeStatusBox(printNodeStatus),
		),
		components.Table(&components.TableProps{
			Columns: components.TableColumns{
				{TitleContents: g.Text("Template")},
				{TitleContents: g.Text("Requirement")},
				{TitleContents: g.Text("PrintNode Job")},
				{TitleContents: g.Text("PDF Title")},
				{TitleContents: g.Text("Created By")},
				{TitleContents: g.Text("Created At")},
				{TitleContents: g.Text("Error")},
				{TitleContents: g.Text("Actions")},
			},
			Rows: rows,
		}),
	)
}

func PDFPrintJobsPage(p PDFPrintJobsPageProps) g.Node {
	return layout.Page(layout.PageProps{
		Title: "Print Jobs",
		Ctx:   p.Ctx,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "printer-settings",
				Title:          "Printing",
				URLPart:        "printing",
			},
		},
		Content: h.Div(
			h.Class("print-jobs-page"),
			printLogsSection(
				p.PrintLogs,
				p.Printers,
				p.Ctx.User.Permissions.Printing.Admin,
				p.PrintNodeStatus,
			),
		),
		AppendHead: []g.Node{components.InlineStyle("/internal/views/pdfview/pdf_page.css")},
	})
}
