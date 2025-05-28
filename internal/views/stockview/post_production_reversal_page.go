package stockview

import (
	"app/internal/components"
	"app/pkg/reqcontext"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
	"github.com/shopspring/decimal"
)

type PostProductionReversalPageProps struct {
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

func PostProductionReversalPage(p *PostGenericPageProps) g.Node {
	p.StockCodePlaceholder = "Enter stock code for production reversal"
	p.QtyPlaceholder = "Enter quantity"

	content := g.Group([]g.Node{
		h.P(
			h.Class("transaction-info"),
			components.Icon(&components.IconProps{
				Identifier: "information-outline",
			}),
			g.Text(
				`Use this utility to post a manual Production Reversal entry from
				STOCK to PRODUCTION accounts for a given location and bin. The
				stock code will also be	issued in  at the given location.`),
			h.Br(),
			g.Text(`NOTE: this utility does not reverse the consumption of stock
				and should be used for corrections with caution.`),
		),

		h.FormEl(
			h.Method("POST"),

			formPartialStockCodeLocBinLot(p),

			components.Button(
				&components.ButtonProps{
					ButtonType: "Primary",
				},
				g.Text("Post Production Reversal Transaction"),
			),
		),
	})
	return postTransactionPageLayout(&postTransactionPageLayoutProps{
		transactionType: "Production Reversal",
		content:         content,
		ctx:             p.Ctx,
		successText:     p.SuccessText,
		errorText:       p.ErrorText,
	})
}
