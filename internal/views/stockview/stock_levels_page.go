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

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

var HomePageDefaultPageSize = 200

type StockLevelsPageProps struct {
	Ctx          reqcontext.ReqContext
	StockLevels  *[]model.StockLevel
	Account      string
	StockCode    string
	Location     string
	Bin          string
	LotNumber    string
	LTETimestamp *time.Time
	Page         int
	PageSize     int
	Total        int
}

func StockLevelsPage(p StockLevelsPageProps) g.Node {

	perms := p.Ctx.User.Permissions

	content := h.FormEl(
		h.Method("GET"),

		h.Nav(
			h.Class("stock-nav"),
			h.A(h.Href("/stock/transactions"), g.Text("See all transactions")),
			g.If(
				perms.SupplyChain.Admin,
				h.A(h.Href("/stock/post-transaction/stock-movement"), g.Text("Post transaction")),
			),
		),

		h.H3(g.Text("Stock Levels")),

		filters(p.Account, p.StockCode, p.Location, p.Bin, p.LotNumber, p.LTETimestamp),

		components.Divider(),

		stockLevelsTable(&stockLevelsTableProps{
			stockLevels: *p.StockLevels,
			page:        p.Page,
			pageSize:    p.PageSize,
			total:       p.Total,
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
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/stockview/stock_table.css"),
		},
	})
}

func filters(
	account, stockCode, location, bin, lotNumber string,
	lteTimestamp *time.Time,
) g.Node {

	lteTimestampStr := ""
	if lteTimestamp != nil {
		lteTimestampStr = lteTimestamp.Format("2006-01-02T15:04")
	}

	return h.Div(
		h.Class("stock-levels-filters"),

		h.Label(
			h.Class("filter"),
			g.Text("Account"),
			h.Select(
				h.Class("lg"),
				h.Name("Account"),
				g.Group(g.Map(model.StockAccounts, func(a model.StockAccount) g.Node {
					return h.Option(
						h.Value(string(a)),
						g.Text(string(a)),
						g.If(account == string(a), h.Selected()),
					)
				})),
			),
		),

		h.Label(
			h.Class("filter"),
			g.Text("Stock Code"),
			h.Input(
				h.Class("lg"),
				h.Name("StockCode"),
				h.Value(stockCode),
				h.AutoComplete("off"),
				h.Placeholder("Enter stock code"),
			),
		),

		h.Label(
			h.Class("filter"),
			g.Text("Location"),
			h.Input(
				h.Class("lg"),
				h.Name("Location"),
				h.Value(location),
				h.AutoComplete("off"),
				h.Placeholder("Enter location"),
			),
		),

		h.Label(
			h.Class("filter"),
			g.Text("Bin"),
			h.Input(
				h.Class("lg"),
				h.Name("Bin"),
				h.Value(bin),
				h.AutoComplete("off"),
				h.Placeholder("Enter bin"),
			),
		),

		h.Label(
			h.Class("filter"),
			g.Text("Lot Number"),
			h.Input(
				h.Class("lg"),
				h.Name("LotNumber"),
				h.Value(lotNumber),
				h.AutoComplete("off"),
				h.Placeholder("Enter lot number"),
			),
		),

		h.Label(
			h.Class("filter"),
			g.Text("Date/Time Limit"),
			h.Input(
				h.Class("lg"),
				h.Type("datetime-local"),
				h.Name("LTETimestamp"),
				h.Value(lteTimestampStr),
				h.AutoComplete("off"),
			),
		),

		h.Div(
			h.Class("go-button-wrapper"),
			components.Button(&components.ButtonProps{
				ButtonType: components.ButtonPrimary,
				Classes:    c.Classes{"go-button": true},
				Size:       components.ButtonLg,
			},
				h.Type("button"),
				g.Attr("onclick", "submitTableForm(this.form)"),
				g.Text("GO"),
			),
		),
	)
}

type stockLevelsTableProps struct {
	stockLevels []model.StockLevel
	page        int
	pageSize    int
	total       int
}

func stockLevelsTable(p *stockLevelsTableProps) g.Node {

	columns := components.TableColumns{{
		TitleContents: g.Text("Account"),
	}, {
		TitleContents: g.Text("Stock Code"),
	}, {
		TitleContents: g.Text("Location"),
	}, {
		TitleContents: g.Text("Bin"),
	}, {
		TitleContents: g.Text("Lot Number"),
	}, {
		TitleContents: g.Text("Stock Level"),
	}, {
		TitleContents: g.Text("Timestamp"),
	}, {
		TitleContents: g.Text(""),
	}}

	var rows components.TableRows

	for _, sl := range p.stockLevels {

		lotNumber := sl.LotNumber
		if lotNumber == "" {
			lotNumber = "\u2013"
		}

		trxParams := url.Values{}
		trxParams.Add("Account", string(sl.Account))
		trxParams.Add("StockCode", sl.StockCode)
		trxParams.Add("Location", sl.Location)
		trxParams.Add("Bin", sl.Bin)
		trxParams.Add("LotNumber", sl.LotNumber)
		trxParams.Add("LTETimestamp", sl.Timestamp.Format("2006-01-02T15:04"))
		transactionsLink := fmt.Sprintf("/stock/transactions?%s", trxParams.Encode())

		rowCells := []components.TableCell{{
			Contents: g.Text(string(sl.Account)),
		}, {
			Contents: components.StockItemAnchor(sl.StockCode),
		}, {
			Contents: g.Text(sl.Location),
		}, {
			Contents: g.Text(sl.Bin),
		}, {
			Contents: g.Text(lotNumber),
		}, {
			Contents:   g.Text(format.DecimalWithCommas(sl.StockLevel.String())),
			Attributes: []g.Node{h.StyleAttr("text-align:right;")},
		}, {
			Contents: h.Span(h.Class("local-datetime"), g.Text(sl.Timestamp.Format(time.RFC3339))),
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
