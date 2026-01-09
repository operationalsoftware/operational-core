package pdfview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/reqcontext"
	"fmt"
	"net/url"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type PrinterAssignmentsPageProps struct {
	Ctx         reqcontext.ReqContext
	Assignments []model.PrintRequirement
}

func PrinterAssignmentsPage(p PrinterAssignmentsPageProps) g.Node {
	canEdit := p.Ctx.User.Permissions.Automation.PrinterAssignmentsEditor

	rows := components.TableRows{}

	for _, pr := range p.Assignments {
		printerName := pr.PrinterName
		if printerName == "" && pr.PrinterID != 0 {
			printerName = fmt.Sprintf("Printer %d", pr.PrinterID)
		}

		cells := []components.TableCell{
			{Contents: g.Text(pr.RequirementName)},
			{Contents: g.Text(printerName)},
		}
		if canEdit {
			cells = append(cells, components.TableCell{
				Contents: h.A(
					h.Class("button primary small"),
					h.Href("/printing/printer-assignments/edit?"+url.Values{"RequirementName": []string{pr.RequirementName}}.Encode()),
					components.Icon(&components.IconProps{Identifier: "pencil"}),
				),
			})
		}

		rows = append(rows, components.TableRow{Cells: cells})
	}

	columns := components.TableColumns{
		{TitleContents: g.Text("Print requirement")},
		{TitleContents: g.Text("Printer Name")},
	}
	if canEdit {
		columns = append(columns, components.TableColumn{TitleContents: g.Text("Edit")})
	}

	assignmentsTable := components.Table(&components.TableProps{Columns: columns, Rows: rows})

	return layout.Page(layout.PageProps{
		Ctx:   p.Ctx,
		Title: "Printer Assignments",
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
			}},
		Content: h.Div(
			h.Class("printer-assignments-page"),
			h.Section(
				h.Class("print-requirements-section"),
				h.H2(g.Text("Print requirements")),
				assignmentsTable,
			),
		),
		AppendHead: []g.Node{components.InlineStyle("/internal/views/pdfview/pdf_page.css")},
	})
}
