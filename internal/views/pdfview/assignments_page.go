package pdfview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/printnode"
	"app/pkg/reqcontext"
	"fmt"
	"net/url"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type PrinterAssignmentsPageProps struct {
	Ctx         reqcontext.ReqContext
	Printers    []printnode.Printer
	Assignments []model.PrintRequirement
	PrintLogs   []model.PDFPrintLog
}

func PrinterAssignmentsPage(p PrinterAssignmentsPageProps) g.Node {
	rows := components.TableRows{}

	for _, pr := range p.Assignments {
		printerName := pr.PrinterName
		if printerName == "" && pr.PrinterID != 0 {
			printerName = fmt.Sprintf("Printer %d", pr.PrinterID)
		}
		rows = append(rows, components.TableRow{
			Cells: []components.TableCell{
				{Contents: g.Text(pr.RequirementName)},
				{Contents: g.Text(printerName)},
				{Contents: h.A(
					h.Class("button secondary"),
					h.Href("/pdf/printer-assignments/edit?"+url.Values{"RequirementName": []string{pr.RequirementName}}.Encode()),
					components.Icon(&components.IconProps{Identifier: "pencil"}),
					g.Text(" Edit"),
				)},
			},
		})
	}

	assignmentsTable := components.Table(&components.TableProps{
		Columns: components.TableColumns{
			{TitleContents: g.Text("Print requirement")},
			{TitleContents: g.Text("Printer")},
			{TitleContents: g.Text("Actions")},
		},
		Rows: rows,
	})

	return layout.Page(layout.PageProps{
		Ctx:   p.Ctx,
		Title: "Printer Assignments",
		Content: h.Div(
			h.Class("printer-assignments-page"),
			h.Section(
				h.Class("print-requirements-section"),
				h.H2(g.Text("Print requirements")),
				assignmentsTable,
			),
			printLogsSection(p.PrintLogs, p.Printers),
		),
		AppendHead: []g.Node{components.InlineStyle("/internal/views/pdfview/pdf_page.css")},
	})
}
