package stockview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type StockDetailPageProps struct {
	Ctx       reqcontext.ReqContext
	StockCode string
}

func StockDetailPage(p StockDetailPageProps) g.Node {

	content := h.FormEl(

		h.H3(g.Text(p.StockCode)),
	)

	return layout.Page(layout.PageProps{
		Title:   "Stock",
		Content: content,
		Ctx:     p.Ctx,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "package-variant-closed",
				Title:          "Stock",
				URLPart:        "stock",
			},
			{
				Title: "Stock Details",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/stockview/index.css"),
		},
	})
}
