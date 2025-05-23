package stockview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
	"app/pkg/format"
	"app/pkg/reqcontext"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

var HomePageDefaultPageSize = 200

type StockLevelsPageProps struct {
	Ctx          reqcontext.ReqContext
	StockLevels  *[]model.StockLevel
	Account      sql.NullString
	StockCode    sql.NullString
	Location     sql.NullString
	Bin          sql.NullString
	LotNumber    sql.NullString
	LTETimestamp sql.NullTime
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
				perms.SupplyChain.Admin || perms.Production.Admin || true,
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
			components.InlineStyle("/internal/views/stockview/index.css"),
		},
	})
}

func filters(
	account, stockCode, location, bin, lotNumber sql.NullString,
	lteTimestamp sql.NullTime,
) g.Node {

	accountStr := "STOCK" // Default
	if account.Valid {
		accountStr = account.String
	}
	lteTimestampStr := ""
	if lteTimestamp.Valid {
		lteTimestampStr = lteTimestamp.Time.Format("2006-01-02T15:04")
	}

	return h.Div(
		h.Class("stock-levels-filters"),

		h.Label(
			h.Class("filter"),
			g.Text("Account"),
			h.Select(
				h.Class("lg"),
				h.Name("Account"),
				g.Group(g.Map(model.StockAccounts, func(a string) g.Node {
					return h.Option(
						h.Value(a),
						g.Text(a),
						g.If(accountStr == a, h.Selected()),
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
				h.Value(stockCode.String),
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
				h.Value(location.String),
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
				h.Value(bin.String),
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
				h.Value(lotNumber.String),
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

		lotNumber := sl.LotNumber.String
		if lotNumber == "" {
			lotNumber = "\u2013"
		}

		trxParams := url.Values{}
		trxParams.Add("Account", sl.Account)
		trxParams.Add("StockCode", sl.StockCode)
		trxParams.Add("Location", sl.Location)
		trxParams.Add("Bin", sl.Bin)
		if sl.LotNumber.Valid {
			trxParams.Add("LotNumber", sl.LotNumber.String)
		}
		trxParams.Add("LTETimestamp", sl.Timestamp.Format("2006-01-02T15:04"))
		transactionsLink := fmt.Sprintf("/stock/transactions?%s", trxParams.Encode())

		rowCells := []components.TableCell{{
			Contents: g.Text(sl.Account),
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
