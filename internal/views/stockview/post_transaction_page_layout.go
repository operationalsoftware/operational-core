package stockview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type postTransactionPageLayoutProps struct {
	ctx             reqcontext.ReqContext
	content         g.Node
	transactionType string
	errorText       string
	successText     string
}

func postTransactionPageLayout(p *postTransactionPageLayoutProps) g.Node {

	type transactionType struct {
		title    string
		linkPart string
	}

	transactionTypes := []transactionType{{
		title:    "Stock Movement",
		linkPart: "stock-movement",
	}, {
		title:    "Production",
		linkPart: "production",
	}, {
		title:    "Production Reversal",
		linkPart: "production-reversal",
	}, {
		title:    "Consumption",
		linkPart: "consumption",
	}, {
		title:    "Consumption Reversal",
		linkPart: "consumption-reversal",
	}, {
		title:    "Stock Adjustment",
		linkPart: "stock-adjustment",
	}}

	content := components.Card(
		h.Nav(
			g.Group(g.Map(transactionTypes, func(tt transactionType) g.Node {
				return h.A(
					g.If(tt.title == p.transactionType, h.Class("active")),
					h.Href("/stock/post-transaction/"+tt.linkPart),
					g.Text(tt.title),
				)
			})),
		),

		components.Divider(),

		h.H3(g.Text("Post "+p.transactionType)),

		p.content,

		g.If(
			p.errorText != "",
			h.Div(
				h.Class("error-msg"),
				g.Text(p.errorText),
			),
		),

		g.If(
			p.successText != "",
			h.Div(
				h.Class("success-msg"),
				g.Text(p.successText),
			),
		),
	)

	return layout.Page(layout.PageProps{
		Title:   "Post " + p.transactionType,
		Content: content,
		Ctx:     p.ctx,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "package-variant-closed",
				Title:          "Stock",
				URLPart:        "stock",
			},
			{
				Title:   "Post Transaction",
				URLPart: "post-transaction",
			},
			{
				Title: p.transactionType,
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle(
				"/internal/views/stockview/post_transactions.css",
			),
		},
	})
}
