package stockview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
	"app/pkg/format"
	"app/pkg/reqcontext"
	"fmt"
	"net/url"
	"time"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

var TransactionsPageDefaultPageSize = 200

type StockTransactionsPageProps struct {
	Ctx               reqcontext.ReqContext
	StockTransactions *[]model.StockTransactionEntry
	Account           string
	StockCode         string
	Location          string
	Bin               string
	LotNumber         string
	LTETimestamp      *time.Time
	Page              int
	PageSize          int
	Total             int
}

func StockTransactionsPage(p *StockTransactionsPageProps) g.Node {

	perms := p.Ctx.User.Permissions

	content := h.FormEl(
		h.Method("GET"),

		h.Nav(
			h.Class("stock-nav"),
			g.If(
				perms.SupplyChain.Admin,
				h.A(h.Href("/stock/post-transaction/stock-movement"), g.Text("Post transaction")),
			),
		),

		h.H3(g.Text("Stock Transactions")),

		filters(p.Account, p.StockCode, p.Location, p.Bin, p.LotNumber, p.LTETimestamp),

		components.Divider(),

		transactionsTable(&transactionsTableProps{
			stockTransactions: *p.StockTransactions,
			page:              p.Page,
			pageSize:          p.PageSize,
			total:             p.Total,
		}),
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
				Title: "Transactions",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/stockview/stock_table.css"),
		},
	})
}

type transactionsTableProps struct {
	stockTransactions []model.StockTransactionEntry
	page              int
	pageSize          int
	total             int
}

func transactionsTable(p *transactionsTableProps) g.Node {

	columns := components.TableColumns{{
		TitleContents: g.Text("ID"),
	}, {
		TitleContents: g.Text("Account"),
	}, {
		TitleContents: g.Text("Transaction Type"),
	}, {
		TitleContents: g.Text("Stock Code"),
	}, {
		TitleContents: g.Text("Location"),
	}, {
		TitleContents: g.Text("Bin"),
	}, {
		TitleContents: g.Text("Lot Number"),
	}, {
		TitleContents: g.Text("Qty Out"),
	}, {
		TitleContents: g.Text("Qty In"),
	}, {
		TitleContents: g.Text("Stock Level"),
	}, {
		TitleContents: g.Text("Timestamp"),
	}, {
		TitleContents: g.Text("By"),
	}, {
		TitleContents: g.Text(""),
	}}

	var rows components.TableRows

	for _, st := range p.stockTransactions {

		lotNumber := st.LotNumber
		if lotNumber == "" {
			lotNumber = "\u2013"
		}

		trxParams := url.Values{}
		trxParams.Add("Account", st.Account)
		trxParams.Add("StockCode", st.StockCode)
		trxParams.Add("Location", st.Location)
		trxParams.Add("Bin", st.Bin)
		trxParams.Add("LotNumber", st.LotNumber)
		trxParams.Add("LTETimestamp", st.Timestamp.Format("2006-01-02T15:04"))
		transactionsLink := fmt.Sprintf("/stock/transactions?%s", trxParams.Encode())

		rowCells := []components.TableCell{{
			Contents: g.Text(fmt.Sprintf("%d", st.StockTransactionEntryID)),
		}, {
			Contents: g.Text(st.Account),
		}, {
			Contents: g.Text(st.TransactionType),
		}, {
			Contents: components.StockItemAnchor(st.StockCode),
		}, {
			Contents: g.Text(st.Location),
		}, {
			Contents: g.Text(st.Bin),
		}, {
			Contents: g.Text(lotNumber),
		}, {
			Contents: g.Group([]g.Node{
				g.If(
					st.Quantity.IsNegative(),
					g.Text(format.DecimalWithCommas(st.Quantity.Abs().String())),
				),
				g.If(
					st.Quantity.IsPositive(),
					g.Text("\u2013"),
				),
			}),
			Attributes: []g.Node{h.StyleAttr("text-align:right;")},
		}, {
			Contents: g.Group([]g.Node{
				g.If(
					st.Quantity.IsPositive(),
					g.Text(format.DecimalWithCommas(st.Quantity.String())),
				),
				g.If(
					st.Quantity.IsNegative(),
					g.Text("\u2013"),
				),
			}),
			Attributes: []g.Node{h.StyleAttr("text-align:right;")},
		}, {
			Contents:   g.Text(format.DecimalWithCommas(st.RunningTotal.String())),
			Attributes: []g.Node{h.StyleAttr("text-align:right;")},
		}, {
			Contents: h.Span(h.Class("local-datetime"), g.Text(st.Timestamp.Format(time.RFC3339))),
		}, {
			Contents: g.Text(st.TransactionByUsername),
		}, {
			Contents: h.A(h.Href(transactionsLink), g.Text("Transactions")),
		}}

		rows = append(rows, components.TableRow{
			Cells:      rowCells,
			Attributes: []g.Node{},
		})
	}

	return components.Table(&components.TableProps{
		Classes: c.Classes{"stock-table": true},
		Columns: columns,
		Rows:    rows,
		Sort:    []appsort.SortItem{},
		Pagination: &components.TablePaginationProps{
			TotalRecords: p.total,
			CurrentPage:  p.page,
			PageSize:     p.pageSize,
		},
	})
}
