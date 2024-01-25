package components

import (
	g "github.com/maragudk/gomponents"
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

type TableProps struct {
	Columns []TableColumn
	Data    []TableRowRenderer
}

func renderHead(columns []TableColumn) g.Node {
	return g.Group(g.Map(columns, func(c TableColumn) g.Node {
		return h.Th(
			h.DataAttr("key", c.Key),
			h.DataAttr("sort", ""),
			h.Class("table-head"),
			h.Span(g.Text(c.Name)),
			g.If(c.Sortable,
				Icon(&IconProps{
					Identifier: "arrow-up-down",
				}),
			),
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
	return h.Div(
		classes,
		h.Table(
			h.THead(
				h.Tr(
					renderHead(p.Columns),
				),
			),
			h.TBody(
				renderRows(p),
			),
		),
		InlineStyle(Assets, "/Table.css"),
		InlineScript(Assets, "/Table.js"),
	)
}
