package pdfview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/reqcontext"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type PDFGenerationLogsPageProps struct {
	Ctx       reqcontext.ReqContext
	Logs      []model.PDFGenerationLog
	Total     int
	Page      int
	PageSize  int
	PageQuery string
	SizeQuery string
}

func PDFGenerationLogsPage(p PDFGenerationLogsPageProps) g.Node {
	rows := components.TableRows{}
	for _, log := range p.Logs {
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

	table := h.Form(
		h.Method("GET"),
		components.Table(&components.TableProps{
			Columns: components.TableColumns{
				{TitleContents: g.Text("Template")},
				{TitleContents: g.Text("PDF title")},
				{TitleContents: g.Text("Input data")},
				{TitleContents: g.Text("Generated at")},
			},
			Rows: rows,
			Pagination: &components.TablePaginationProps{
				TotalRecords:        p.Total,
				PageSize:            p.PageSize,
				CurrentPage:         p.Page,
				CurrentPageQueryKey: p.PageQuery,
				PageSizeQueryKey:    p.SizeQuery,
			},
		}),
	)

	return layout.Page(layout.PageProps{
		Title: "PDF Generation Logs",
		Ctx:   p.Ctx,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "text-box-outline",
				Title:          "PDFs",
				URLPart:        "pdf",
			},
			{
				Title:   "PDF Generation Logs",
				URLPart: "logs",
			},
		},
		Content: h.Section(
			h.Class("pdf-generation-log-section"),
			h.H2(g.Text("PDF generation logs")),
			table,
		),
		AppendHead: []g.Node{components.InlineStyle("/internal/views/pdfview/pdf_page.css")},
	})
}
