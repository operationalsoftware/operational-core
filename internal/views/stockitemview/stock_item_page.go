package stockitemview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/appsort"
	"app/pkg/reqcontext"
	"fmt"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type StockItemPageProps struct {
	Id                int
	Ctx               reqcontext.ReqContext
	StockItem         model.StockItem
	QRCode            string
	GalleryImageURLs  []string
	GalleryURL        string
	StockItemChanges  []model.StockItemChange
	StockItemComments []model.Comment
	Sort              appsort.Sort
	Page              int
	PageSize          int
}

func StockItemPage(p *StockItemPageProps) g.Node {

	content := g.Group([]g.Node{

		h.Div(
			h.Class("header"),

			h.H3(g.Text(p.StockItem.StockCode)),

			h.Div(
				h.Class("actions"),

				h.A(
					h.Class("button primary"),
					h.Href(fmt.Sprintf("/stock-items/%d/edit", p.StockItem.StockItemID)),
					components.Icon(&components.IconProps{
						Identifier: "pencil",
					}),
				),
			),
		),

		h.Div(
			h.Class("two-column-flex"),

			stockItemProperties(p.StockItem),

			h.Div(
				h.Class("gallery-container"),

				components.Gallery(p.GalleryImageURLs),

				h.A(
					h.Class("button primary"),
					h.Href(p.GalleryURL),

					g.Text("Gallery"),

					components.Icon(&components.IconProps{
						Identifier: "arrow-right-thin",
					}),
				),
			),
		),

		h.Div(
			h.Class("two-column-flex"),

			components.CommentsThread(&components.CommentsThreadProps{
				Comments: p.StockItemComments,
				Entity:   "StockItem",
				EntityID: p.StockItem.StockItemID,
			}),

			stockItemChangeLog(p.StockItemChanges),
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
			{Title: p.StockItem.StockCode},
		},
		Content: content,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/stockitemview/stock_item_page.css"),
		},
	})
}

func stockItemProperties(si model.StockItem) g.Node {
	return h.Div(
		h.Class("properties"),

		g.Map([]struct {
			label string
			value string
		}{
			{"Stock Code", si.StockCode},
			{"Description", si.Description},
		}, func(i struct {
			label string
			value string
		}) g.Node {
			return g.Group([]g.Node{
				h.Div(h.Strong(g.Text(i.label))),
				h.Div(g.Text(i.value)),
				components.CopyButton(i.value),
			})
		}),
	)
}

var changelogFieldDefs = []components.ChangelogProperty{
	{FieldKey: "StockCode", Label: g.Text("Stock Code")},
	{FieldKey: "Description", Label: g.Text("Description")},
}

func stockItemChangeLog(changes []model.StockItemChange) g.Node {

	var changelogEntries []components.ChangelogEntry
	for _, change := range changes {
		entry := components.ChangelogEntry{
			ChangedAt:        change.ChangedAt,
			ChangeByUsername: change.ChangeByUsername,
			IsCreation:       change.IsCreation,
			Changes: map[string]any{
				"StockCode":   change.StockCode,
				"Description": change.Description,
			},
		}
		changelogEntries = append(changelogEntries, entry)
	}

	return components.Changelog(changelogEntries, changelogFieldDefs)
}
