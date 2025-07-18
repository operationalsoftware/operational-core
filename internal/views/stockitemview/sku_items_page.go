package stockitemview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/reqcontext"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type SKUItemsPageProps struct {
	Ctx      reqcontext.ReqContext
	SKUItems model.SKUConfigData
}

func SKUItemsPage(p *SKUItemsPageProps) g.Node {

	content := g.Group([]g.Node{
		h.Div(
			h.Class("button-container"),
			components.Button(&components.ButtonProps{
				ButtonType: "primary",
				Link:       "/stock-items/sku-config/add",
				Classes: c.Classes{
					"add-user-btn": true,
				},
			},
				components.Icon(&components.IconProps{
					Identifier: "plus",
				}),
				g.Text("New Item"),
			),
		),

		h.Div(
			h.Class("sku-lists-container"),

			g.If(
				p.Ctx.Req.URL.Query().Get("toast") == "deleted",
				components.Toast(&components.ToastProps{
					Contents: g.Text("Item deleted successfully!"),
					Type:     "success",
				}),
			),

			skuList(skuListProps{
				title:    "Product Type",
				skuField: "ProductType",
				items:    p.SKUItems.ProductType,
			}),
			skuList(skuListProps{
				title:    "Yarn Type",
				skuField: "YarnType",
				items:    p.SKUItems.YarnType,
			}),
			skuList(skuListProps{
				title:    "Style Number",
				skuField: "StyleNumber",
				items:    p.SKUItems.StyleNumber,
			}),
			skuList(skuListProps{
				title:    "Colour",
				skuField: "Colour",
				items:    p.SKUItems.Colour,
			}),
			skuList(skuListProps{
				title:    "Toe Closing",
				skuField: "ToeClosing",
				items:    p.SKUItems.ToeClosing,
			}),
			skuList(skuListProps{
				title:    "Size",
				skuField: "Size",
				items:    p.SKUItems.Size,
			}),
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
				URLPart:        "stock-items",
			},
			{
				IconIdentifier: "barcode-scan",
				Title:          "SKU Configuration",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/stockitemview/sku_items_page.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/stockitemview/sku_items_page.js"),
		},
	})
}

type skuListProps struct {
	title    string
	skuField string
	items    []model.SKUConfig
}

// func skuList(title, skuField string, items []model.SKUConfig) g.Node {
func skuList(p skuListProps) g.Node {
	return h.Div(
		h.Class("sku-list"),

		h.H3(
			h.Class("sku-list-title"),

			g.Text(p.title),
		),
		h.Input(
			h.Type("text"),
			h.Class("sku-search-input"),
			h.Placeholder("Search..."),
			g.Attr("onkeyup", "filterList(event)"),
		),
		h.Div(
			h.Class("sku-list-items"),

			g.Group(g.Map(p.items, func(item model.SKUConfig) g.Node {
				return h.Div(
					h.Class("sku-list-item"),
					g.Attr("data-field", p.skuField),
					g.Attr("data-label", item.Label),
					g.Attr("data-code", item.Code),
					h.Span(
						h.Class("sku-item-label"),
						g.Text(item.Code+" \u2013 "+item.Label),
					),
					h.Button(
						h.Type("button"),
						h.Class("sku-delete-btn"),
						g.Attr("onclick", "deleteSKUItem(event)"),

						h.Span(
							h.Class("icon"),
							components.Icon(&components.IconProps{
								Identifier: "close",
							}),
						),
					),
				)
			})),
		),
	)
}
