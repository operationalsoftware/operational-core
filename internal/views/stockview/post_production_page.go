package stockview

import (
	"app/internal/components"
	"app/pkg/reqcontext"
	"database/sql"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
	"github.com/shopspring/decimal"
)

type PostProductionPageProps struct {
	Ctx         reqcontext.ReqContext
	SuccessText string
	ErrorText   string

	StockCode string
	Location  string
	Bin       string
	LotNumber sql.NullString
	Qty       decimal.Decimal
}

func PostProductionPage(p *PostProductionPageProps) g.Node {

	content := g.Group([]g.Node{
		h.P(
			h.Class("transaction-info"),
			components.Icon(&components.IconProps{
				Identifier: "information-outline",
			}),
			g.Text(
				`Use this utility to post a manual Production entry from
				PRODUCTION to STOCK accounts for a given location and bin. The
				stock code will also be	received in SyteLine at the given location.`),
			h.Br(),
			g.Text(`NOTE: this utility does not consume stock and should be used for
				corrections with caution.`),
		),

		h.FormEl(
			h.Method("POST"),

			h.Div(
				h.Class("form-row"),

				h.Label(
					g.Text("Stock Code"),
					h.Input(
						h.Type("text"),
						h.Name("StockCode"),
						h.Value(p.StockCode),
						h.Placeholder("Enter stock code to produce"),
						h.AutoComplete("off"),
					),
				),
			),

			h.Div(
				h.Class("form-row"),

				h.Label(
					g.Text("Location"),
					h.Input(
						h.Type("text"),
						h.Name("Location"),
						h.Value(p.Location),
						h.Placeholder("Enter location"),
						h.AutoComplete("off"),
					),
				),
			),

			h.Div(
				h.Class("form-row"),

				h.Label(
					g.Text("Bin"),
					h.Input(
						h.Type("text"),
						h.Name("Bin"),
						h.Value(p.Bin),
						h.Placeholder("Enter bin"),
						h.AutoComplete("off"),
					),
				),
			),

			h.Div(
				h.Class("form-row"),

				h.Label(
					g.Text("Lot Number (only if lot tracked)"),
					h.Input(
						h.Type("text"),
						h.Name("LotNumber"),
						h.Value(p.LotNumber.String),
						h.Placeholder("Enter lot number"),
						h.AutoComplete("off"),
					),
				),
			),

			h.Div(
				h.Class("form-row"),

				h.Label(
					g.Text("Qty"),
					h.Input(
						h.Type("number"),
						h.Min("0"),
						h.Name("Qty"),
						h.Step("any"),
						g.If(
							p.Qty.GreaterThan(decimal.Zero),
							h.Value(p.Qty.String()),
						),
						h.Placeholder("Enter qty to produce"),
						h.AutoComplete("off"),
					),
				),
			),

			components.Button(
				&components.ButtonProps{
					ButtonType: "Primary",
				},
				g.Text("Post Production Transaction"),
			),
		),
	})
	return postTransactionPageLayout(&postTransactionPageLayoutProps{
		transactionType: "Production",
		content:         content,
		ctx:             p.Ctx,
		successText:     p.SuccessText,
		errorText:       p.ErrorText,
	})
}
