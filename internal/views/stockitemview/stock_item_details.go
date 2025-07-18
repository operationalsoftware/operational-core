package stockitemview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
	"app/pkg/nilsafe"
	"app/pkg/reqcontext"
	"fmt"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type StockItemDetailsPageProps struct {
	Id               int
	Ctx              reqcontext.ReqContext
	StockItem        *model.StockItem
	QRCode           string
	StockItemChanges []model.StockItemChange
	Sort             appsort.Sort
	Page             int
	PageSize         int
}

func StockItemDetailsPage(p *StockItemDetailsPageProps) g.Node {

	stockItem := p.StockItem

	var tableRows []g.Node
	for _, u := range p.StockItemChanges {

		tableRows = append(tableRows,
			h.Li(
				h.Strong(g.Text(fmt.Sprintf("Changed by %s at %s", u.ChangeByUsername, u.ChangeAt.Format("2006-01-02 15:04:05")))),
				h.Ul(
					g.If(u.StockCodeHistory != nil,
						h.Li(
							g.Text("Stock Code: "+nilsafe.Str(u.StockCodeHistory))),
					),
					g.If(u.Description != nil,
						h.Li(
							g.Text("Description: "+nilsafe.Str(u.Description))),
					),
				),
			),
		)
	}

	stockCode := stockItem.StockCode
	stockDescription := stockItem.Description

	content := g.Group([]g.Node{

		h.Div(
			h.Class("button-container"),
			components.Button(&components.ButtonProps{
				ButtonType: "primary",
				Classes: c.Classes{
					"edit-button": true,
				},
				Link: fmt.Sprintf("/stock-items/%s/edit", p.StockItem.StockCode),
			},
				components.Icon(&components.IconProps{
					Identifier: "pencil",
				}),
			),
		),
		h.Div(
			h.H3(g.Text("Stock Item "+p.StockItem.StockCode)),

			h.Div(
				h.Class("flex"),

				h.Div(
					h.Class("properties-grid"),

					g.Group([]g.Node{
						h.Span(
							h.Strong(g.Text("Stock Code")),
						),
						h.Div(
							h.Class("form-value"),

							h.Span(
								g.Text(stockCode),
							),
							components.CopyButton(components.CopyButtonProps{
								TextToCopy: stockCode,
								ButtonID:   "stock-copy",
							}),
						),
						h.Span(
							h.Strong(g.Text("Description")),
						),
						h.Div(
							h.Class("form-value"),

							h.Span(
								g.Text(stockDescription),
							),
							components.CopyButton(components.CopyButtonProps{
								TextToCopy: stockDescription,
								ButtonID:   "description-copy",
							}),
						),
					}),
				),
			),
		),

		h.Br(),
		h.Br(),
		h.H3(
			h.Class("changes-heading"),
			g.Text("Changelog"),
		),

		h.Div(
			h.ID("stock-items-changes"),
			g.Group(tableRows),
		),
	})

	return layout.Page(layout.PageProps{
		Title: "Stock Item Details",
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "package-variant-closed",
				Title:          "Stock Items",
				URLPart:        "stock-items",
			},
			{Title: stockItem.StockCode},
		},
		Content: content,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/userview/user_page.css"),
			components.InlineStyle("/internal/views/stockitemview/stock_item_details.css"),
		},
	})
}
