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
	MyFilter        string
}

func StockItemsPage(p *StockItemsPageProps) g.Node {
	var columns = components.TableColumns{
		{
			TitleContents: g.Text("Stock Code"),
			SortKey:       "Stock_Code",
		},
		{
			TitleContents: g.Text("Description"),
			SortKey:       "Description",
		},
		{
			TitleContents: g.Text("Created"),
			SortKey:       "Created_At",
		},
	}

	var tableRows components.TableRows
	for _, u := range p.StockItems {

		tableRows = append(tableRows, components.TableRow{
			Cells: []components.TableCell{
				{
					Contents: h.A(
						g.Text(u.StockCode),
						g.Attr("href",
							fmt.Sprintf("/stock-items/%d", u.StockItemID))),
				},
				{
					Contents: g.Group([]g.Node{
						g.If(u.Description == "", g.Text("\u2013")),
						g.If(u.Description != "", g.Text(u.Description)),
					}),
				},
				{
					Contents: g.Text(u.CreatedAt.Format("2006-01-02 15:04:05")),
				},
			},
		})
	}

	content := g.Group([]g.Node{

		h.Div(
			h.Class("button-container"),
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

		// form container for table interaction
		h.FormEl(
			h.ID("stock-items-table-form"),
			g.Attr("method", "GET"),

			components.Table(&components.TableProps{
				Columns: columns,
				Sort:    p.Sort,
				Rows:    tableRows,
				Pagination: &components.TablePaginationProps{
					TotalRecords:        p.StockItemsCount,
					PageSize:            p.PageSize,
					CurrentPage:         p.Page,
					CurrentPageQueryKey: "Page",
					PageSizeQueryKey:    "PageSize",
				},
			},
				h.ID("stock-items-table"),
			),
		),
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
