package stockview

import (
	"app/internal/components"
	"app/internal/model"
	"app/pkg/reqcontext"
	"fmt"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
	"github.com/shopspring/decimal"
)

type PostGenericPageProps struct {
	Ctx         reqcontext.ReqContext
	SuccessText string
	ErrorText   string

	StockItemID     int
	Location        string
	Bin             string
	LotNumber       string
	Qty             decimal.Decimal
	TransactionNote string

	IsStockAdjustment bool

	// Placeholders
	StockCodePlaceholder string
	QtyPlaceholder       string

	StockItems []model.StockItem
}

func PostProductionPage(p *PostGenericPageProps) g.Node {
	p.StockCodePlaceholder = "Enter stock code to produce"
	p.QtyPlaceholder = "Enter quantity to produce"

	content := g.Group([]g.Node{
		h.P(
			h.Class("transaction-info"),
			components.Icon(&components.IconProps{
				Identifier: "information-outline",
			}),
			g.Text(
				`Use this utility to post a manual Production entry from
				PRODUCTION to STOCK accounts for a given location and bin. The
				stock code will also be	received in  at the given location.`),
			h.Br(),
			g.Text(`NOTE: this utility does not consume stock and should be used for
				corrections with caution.`),
		),

		h.FormEl(
			h.Method("POST"),

			formPartialStockCodeLocBinLot(p),

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

func formPartialStockCodeLocBinLot(p *PostGenericPageProps) g.Node {

	selectedStockItem := ""
	if p.StockItemID != 0 {
		selectedStockItem = fmt.Sprintf("%d", p.StockItemID)
	}

	return g.Group([]g.Node{
		h.Div(
			h.Class("form-row"),

			h.Label(
				g.Text("Stock Code"),

				components.SearchSelect(&components.SearchSelectProps{
					Name:                 "StockItemID",
					Placeholder:          "Select Stock Code",
					Mode:                 "single",
					Options:              MapStockItemsToOptions(p.StockItems, selectedStockItem),
					Selected:             selectedStockItem,
					OptionsEndpoint:      "/get-stock-codes",
					SearchQueryParamName: "SearchText",
				}),
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
					h.Value(p.LotNumber),
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
					g.If(!p.IsStockAdjustment, h.Min("0")),
					h.Name("Qty"),
					h.Step("any"),
					g.If(
						p.Qty.GreaterThan(decimal.Zero) || p.Qty.LessThan(decimal.Zero),
						h.Value(p.Qty.String()),
					),
					h.Placeholder(p.QtyPlaceholder),
					h.AutoComplete("off"),
				),
			),
		),

		h.Div(
			h.Class("form-row"),

			h.Label(
				g.Text("Note (optional)"),
				h.Textarea(
					h.Name("TransactionNote"),
					h.Value(p.TransactionNote),
					h.Placeholder("Enter transaction note"),
					h.AutoComplete("off"),
				),
			),
		),
	})
}
