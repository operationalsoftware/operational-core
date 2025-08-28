package stockview

import (
	"app/internal/components"
	"app/pkg/reqcontext"

	"github.com/shopspring/decimal"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type PostConsumptionPageProps struct {
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

func PostConsumptionPage(p *PostGenericPageProps) g.Node {
	p.StockCodePlaceholder = "Enter stock code for consumption"
	p.QtyPlaceholder = "Enter quantity to consume"

	content := g.Group([]g.Node{
		h.P(
			h.Class("transaction-info"),
			components.Icon(&components.IconProps{
				Identifier: "information-outline",
			}),
			g.Text(
				`Use this utility to post a manual Consumption entry from
				STOCK to CONSUMED accounts for a given location and bin. The
				stock code will also be	issued in  at the given location.`),
		),
		h.P(
			h.Class("transaction-info"),
			g.Text(`NOTE: this utility should be used for corrections with caution.`),
		),

		h.Form(
			h.Method("POST"),

			formPartialStockCodeLocBinLot(p),

			components.Button(
				&components.ButtonProps{
					ButtonType: "Primary",
				},
				g.Text("Post Consumption Transaction"),
			),
		),
	})
	return postTransactionPageLayout(&postTransactionPageLayoutProps{
		transactionType: "Consumption",
		content:         content,
		ctx:             p.Ctx,
		successText:     p.SuccessText,
		errorText:       p.ErrorText,
	})
}
