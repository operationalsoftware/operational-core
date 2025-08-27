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

type StockItemDetailsPageProps struct {
	Id                int
	Ctx               reqcontext.ReqContext
	StockItem         *model.StockItem
	QRCode            string
	StockItemChanges  []model.StockItemChange
	StockItemComments []components.Comment
	Sort              appsort.Sort
	Page              int
	PageSize          int
}

var changelogFieldDefs = []components.ChangelogProperty{
	{FieldKey: "StockCode", Label: g.Text("Stock Code")},
	{FieldKey: "Description", Label: g.Text("Description")},
}

func StockItemDetailsPage(p *StockItemDetailsPageProps) g.Node {

	stockItem := p.StockItem

	var changelogEntries []components.ChangelogEntry
	for _, change := range p.StockItemChanges {
		entry := components.ChangelogEntry{
			ChangedAt:        change.ChangedAt,
			ChangeByUsername: change.ChangeByUsername,
			IsCreation:       change.IsCreation,
			Changes: map[string]interface{}{
				"StockCode":   change.StockCode,
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
				Link: fmt.Sprintf("/stock-items/%d/edit", p.StockItem.StockItemID),
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

		h.Div(
			h.Class("history-section"),

			components.CommentsThread(&components.CommentsThreadProps{
				Comments: p.StockItemComments,
				Entity:   "stock item",
				EntityID: p.StockItem.StockItemID,
			}),

			h.Br(),
			h.Br(),

			h.Div(
				h.Class("change-log"),
				components.Changelog(changelogEntries, changelogFieldDefs),
			),
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
			components.InlineStyle("/internal/views/stockitemview/stock_item_page.css"),
		},
	})
}
