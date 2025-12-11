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
		documentCell := g.Text("-")
		linkTitle := log.PDFTitle
		if linkTitle == "" {
			linkTitle = log.TemplateName
		}
		viewerURL := pdfViewerURL(log.FileURL, log.FileID, linkTitle)
		if viewerURL != "" {
			documentCell = h.A(
				h.Href(viewerURL),
				h.Target("_blank"),
				h.Rel("noopener noreferrer"),
				g.Text(linkTitle),
			)
		}

		rows = append(rows, components.TableRow{
			Cells: []components.TableCell{
				{Contents: g.Text(log.TemplateName)},
				{Contents: documentCell},
				{Contents: h.Pre(g.Text(inputPreview))},
				{Contents: g.Text(log.CreatedAt.Format("2006-01-02 15:04:05"))},
			},
		})
	}

	table := components.Table(&components.TableProps{
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
	})

	return layout.Page(layout.PageProps{
		Title: "PDF Generation Logs",
		Ctx:   p.Ctx,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				Title: "PDF Generation Logs",
			},
		},
		Content: h.Section(
			h.Class("pdf-generation-log-section"),
			h.H2(g.Text("PDF generation logs")),
			g.If(len(p.Logs) == 0,
				h.P(g.Text("No PDF generations logged yet.")),
			),
			g.If(len(p.Logs) > 0, table),
		),
		AppendHead: []g.Node{components.InlineStyle("/internal/views/pdfview/pdf_page.css")},
	})
}
