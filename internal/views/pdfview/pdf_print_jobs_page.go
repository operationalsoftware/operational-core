package pdfview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/printnode"
	"app/pkg/reqcontext"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type PDFPrintJobsPageProps struct {
	Ctx             reqcontext.ReqContext
	PrintLogs       []model.PDFPrintLog
	Printers        []printnode.Printer
	PrintNodeStatus printnode.Status
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
				printNodeStatusBox(p.PrintNodeStatus),
			),
		),
		AppendHead: []g.Node{components.InlineStyle("/internal/views/pdfview/pdf_page.css")},
	})
}
