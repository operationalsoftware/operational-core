package stockview

import (
	"app/internal/components"
	"app/pkg/nilsafe"
	"app/pkg/reqcontext"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
	"github.com/shopspring/decimal"
)

type PostStockMovementPageProps struct {
	Ctx         reqcontext.ReqContext
	SuccessText string
	ErrorText   string
	ReturnTo    *string

	TransactionType string
	StockCode       string
	LotNumber       string
	Qty             decimal.Decimal
	FromLocation    string
	FromBin         string
	ToLocation      string
	ToBin           string
	TransactionNote string
}

func PostStockMovementPage(p *PostStockMovementPageProps) g.Node {

	content := h.FormEl(
		h.Method("POST"),

		h.Div(
			h.Class("form-row"),

			h.Label(
				g.Text("Stock Code"),
				h.Input(
					h.Type("text"),
					h.Name("StockCode"),
					h.Value(p.StockCode),
					h.Placeholder("Enter stock code"),
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
					h.Name("Qty"),
					h.Type("number"),
					h.Min("0"),
					h.Step("any"),
					g.If(
						p.Qty.GreaterThan(decimal.Zero),
						h.Value(p.Qty.String()),
					),
					h.Placeholder("Enter quantity"),
					h.AutoComplete("off"),
				),
			),
		),

		h.Div(
			h.Class("form-row"),

			h.Label(
				g.Text("From Location"),
				h.Input(
					h.Type("text"),
					h.Name("FromLocation"),
					h.Value(p.FromLocation),
					h.Placeholder("Enter location (from)"),
					h.AutoComplete("off"),
				),
			),
			h.Label(
				g.Text("To Location"),
				h.Input(
					h.Type("text"),
					h.Name("ToLocation"),
					h.Value(p.ToLocation),
					h.Placeholder("Enter location (to)"),
					h.AutoComplete("off"),
				),
			),
		),

		h.Div(
			h.Class("form-row"),

			h.Label(
				g.Text("From Bin"),
				h.Input(
					h.Type("text"),
					h.Name("FromBin"),
					h.Value(p.FromBin),
					h.Placeholder("Enter bin (from)"),
					h.AutoComplete("off"),
				),
			),
			h.Label(
				g.Text("To Bin"),
				h.Input(
					h.Type("text"),
					h.Name("ToBin"),
					h.Value(p.ToBin),
					h.Placeholder("Enter bin (to)"),
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

		// hidden input to store returnTo
		g.If(
			p.ReturnTo != nil,
			h.Input(
				h.Type("hidden"),
				h.Name("ReturnTo"),
				h.Value(nilsafe.Str(p.ReturnTo)),
			),
		),

		components.Button(
			&components.ButtonProps{
				ButtonType: "Primary",
			},
			g.Text("Post Movement"),
		),
	)

	return postTransactionPageLayout(&postTransactionPageLayoutProps{
		transactionType: p.TransactionType,
		content:         content,
		ctx:             p.Ctx,
		successText:     p.SuccessText,
		errorText:       p.ErrorText,
	})
}
