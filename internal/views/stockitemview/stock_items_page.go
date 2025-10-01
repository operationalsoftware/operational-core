package stockitemview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
	"app/pkg/reqcontext"
	"fmt"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type StockItemsPageProps struct {
	Ctx             reqcontext.ReqContext
	StockItems      []model.StockItem
	StockItemsCount int
	Sort            appsort.Sort
	Page            int
	PageSize        int
}

func StockItemsPage(p *StockItemsPageProps) g.Node {

	isUserStockAdmin := p.Ctx.User.Permissions.Stock.Admin

	content := g.Group([]g.Node{
		stockItemsActions(isUserStockAdmin),

		stockItemsTable(stockItemsTableProps{
			stockItems:      p.StockItems,
			stockItemsCount: p.StockItemsCount,
			sort:            p.Sort,
			page:            p.Page,
			pageSize:        p.PageSize,
		}),
	})

	return layout.Page(layout.PageProps{
		Ctx:     p.Ctx,
		Title:   "Stock Items",
		Content: content,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "package-variant-closed",
				Title:          "Stock Items",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/stockitemview/stock_items_page.css"),
		},
	})
}

func stockItemsActions(isUserStockAdmin bool) g.Node {
	return h.Div(
		h.Class("button-container"),
		g.If(isUserStockAdmin,
			components.Button(&components.ButtonProps{
				ButtonType: "primary",
				Link:       "/stock-items/add",
				Classes: c.Classes{
					"add-user-btn": true,
				},
			},
				components.Icon(&components.IconProps{
					Identifier: "plus",
				}),
				g.Text("Stock Item"),
			),
		),
	)
}

type stockItemsTableProps struct {
	stockItems      []model.StockItem
	stockItemsCount int
	sort            appsort.Sort
	page            int
	pageSize        int
}

func stockItemsTable(p stockItemsTableProps) g.Node {

	var columns = components.TableColumns{
		{TitleContents: g.Text("Stock Code"), SortKey: "StockCode"},
		{TitleContents: g.Text("Description"), SortKey: "Description"},
		{TitleContents: g.Text("Created"), SortKey: "CreatedAt"},
	}

	var tableRows components.TableRows
	for _, si := range p.stockItems {

		stockItemHref := fmt.Sprintf("/stock-items/%d", si.StockItemID)
		description := "\u2013"
		if si.Description != "" {
			description = si.Description
		}

		tableRows = append(tableRows, components.TableRow{
			Cells: []components.TableCell{
				{Contents: h.A(h.Href(stockItemHref), g.Text(si.StockCode))},
				{Contents: g.Text(description)},
				{Contents: g.Text(si.CreatedAt.Format("2006-01-02 15:04:05"))},
			},
			HREF: stockItemHref,
		})
	}

	// form container for table interaction
	return h.Form(
		g.Attr("method", "GET"),

		components.Table(&components.TableProps{
			Columns: columns,
			Sort:    p.sort,
			Rows:    tableRows,
			Pagination: &components.TablePaginationProps{
				TotalRecords:        p.stockItemsCount,
				PageSize:            p.pageSize,
				CurrentPage:         p.page,
				CurrentPageQueryKey: "Page",
				PageSizeQueryKey:    "PageSize",
			},
		}),
	)

}
