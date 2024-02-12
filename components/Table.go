package components

import (
	g "github.com/maragudk/gomponents"
	ghtmx "github.com/maragudk/gomponents-htmx"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type RenderedCell struct {
	Content    g.Node
	Classes    c.Classes
	Attributes []g.Node
}

type TableRowRenderer interface {
	Render() map[string]RenderedCell
}

type TableColumn struct {
	Name     string
	Key      string
	Sortable bool
}

type SortItem struct {
	Key  string
	Sort string
}

type TableProps struct {
	Columns []TableColumn
	Data    []TableRowRenderer
	HxGet   string
	Sort    []SortItem
}

func renderHead(columns []TableColumn, hxGet string, sort []SortItem) g.Node {
	return g.Group(g.Map(columns, func(c TableColumn) g.Node {
		sortDirection := ""
		if c.Sortable {
			for _, s := range sort {
				if s.Key == c.Key {
					sortDirection = s.Sort
				}
			}
		}
		return h.Th(
			h.DataAttr("key", c.Key),
			g.If(sortDirection != "", h.DataAttr("sort", sortDirection)),
			h.Class("table-head"),
			h.Span(g.Text(c.Name)),
			g.Group([]g.Node{
				g.If(sortDirection == "asc", Icon(&IconProps{
					Identifier: "arrow-down",
				})),
				g.If(sortDirection == "desc", Icon(&IconProps{
					Identifier: "arrow-up",
				})),
				g.If(sortDirection == "", Icon(&IconProps{
					Identifier: "arrow-up-down",
				})),
			}),
			g.If(hxGet != "", g.Group([]g.Node{
				ghtmx.Get(hxGet + "?sort=Username-ASC"), ghtmx.Target("closest table"),
				ghtmx.Swap("outerHTML"),
			})),
		)
	}))
}

func renderRows(p *TableProps) g.Node {
	return g.Group(g.Map(p.Data, func(d TableRowRenderer) g.Node {
		return h.Tr(
			h.Class("table-row"),
			g.Group(g.Map(p.Columns, func(col TableColumn) g.Node {
				renderKey := col.Key
				if renderKey == "" {
					panic("TableColumn.Key must be set")
				}

				renderedCell := d.Render()[renderKey]

				if renderedCell.Classes == nil {
					renderedCell.Classes = c.Classes{}
				}

				renderedCell.Classes["table-cell"] = true

				return h.Td(
					renderedCell.Classes,
					g.Group(renderedCell.Attributes),
					renderedCell.Content,
				)
			})),
		)
	}))
}

func Table(p *TableProps) g.Node {
	classes := c.Classes{
		"table-container": true,
	}

	if p.Sort == nil {
		p.Sort = []SortItem{}
	}

	return h.Div(
		classes,
		h.Table(
			ghtmx.PushURL("true"),
			h.THead(
				h.Tr(
					renderHead(p.Columns, p.HxGet, p.Sort),
				),
			),
			h.TBody(
				renderRows(p),
			),
		),
		InlineStyle("/components/Table.css"),
		InlineScript("/components/Table.js"),
	)
}
