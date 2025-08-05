package stockitemview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
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

var changelogFieldDefs = []components.ChangelogFieldDefinition{
	{Name: "StockCode", Label: "Stock Code"},
	{Name: "Description", Label: "Description"},
}

func StockItemDetailsPage(p *StockItemDetailsPageProps) g.Node {

	stockItem := p.StockItem

	var changelogEntries []components.ChangelogEntry
	for _, change := range p.StockItemChanges {
		entry := components.ChangelogEntry{
			ChangedAt:         change.ChangedAt,
			ChangedByUsername: change.ChangeByUsername,
			IsCreation:        change.IsCreation,
			Changes: map[string]interface{}{
				"StockCode":   change.StockCodeHistory,
				"Description": change.Description,
			},
		}
		changelogEntries = append(changelogEntries, entry)
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
							components.CopyButton(stockCode),
						),
						h.Span(
							h.Strong(g.Text("Description")),
						),
						h.Div(
							h.Class("form-value"),

							h.Span(
								g.Text(stockDescription),
							),
							components.CopyButton(stockDescription),
						),
					}),
				),
			),
		),

		h.Br(),
		h.Br(),
		h.Hr(),

		components.Changelog(changelogEntries, changelogFieldDefs),
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
