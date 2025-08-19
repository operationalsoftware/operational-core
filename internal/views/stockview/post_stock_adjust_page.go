package stockview

import (
	"app/internal/components"
	"app/pkg/reqcontext"

	"github.com/shopspring/decimal"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type PostStockAdjustPageProps struct {
	Ctx         reqcontext.ReqContext
	SuccessText string
	ErrorText   string

	StockCode       string
	Location        string
	Bin             string
	LotNumber       string
	Qty             decimal.Decimal
	TransactionNote string
}

func PostStockAdjustPage(p *PostGenericPageProps) g.Node {
	p.StockCodePlaceholder = "Enter stock code for stock adjust"
	p.QtyPlaceholder = "Enter quantity"

	content := g.Group([]g.Node{
		h.P(
			h.Class("transaction-info"),
			components.Icon(&components.IconProps{
				Identifier: "information-outline",
			}),
			g.Text(
				`Use this utility to post a manual Stock Adjustment entry from
				ADJUST to STOCK accounts for a given location and bin. The
				stock code will also be	received in  at the given location.`),
			h.Br(),
			g.Text(`NOTE: this utility should be used for corrections with caution.`),
		),

		h.FormEl(
			h.Method("POST"),

			formPartialStockCodeLocBinLot(p),

			components.Button(
				&components.ButtonProps{
					ButtonType: "Primary",
				},
				g.Text("Post Stock Adjustment Transaction"),
			),
		),
	})
	return postTransactionPageLayout(&postTransactionPageLayoutProps{
		transactionType: "Stock Adjustment",
		content:         content,
		ctx:             p.Ctx,
		successText:     p.SuccessText,
		errorText:       p.ErrorText,
	})
}
