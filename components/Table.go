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
	Name string
	Key  string
}

type TableProps struct {
	Columns []TableColumn
	Data    []TableRowRenderer
}

func renderColumns(columns []TableColumn) g.Node {
	return g.Group(g.Map(columns, func(c TableColumn) g.Node {
		return h.Th(
			h.Class("table-head"),
			g.Text(c.Name),
		)
	}))
}

func renderRows(p *TableProps) g.Node {
	return g.Group(g.Map(p.Data, func(d TableRowRenderer) g.Node {
		return h.Tr(
			h.Class("table-row"),
			g.Group(g.Map(p.Columns, func(c TableColumn) g.Node {
				renderKey := c.Key
				if renderKey == "" {
					renderKey = c.Key
				}

				renderedCell := d.Render()[renderKey]

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
					renderColumns(p.Columns),
				),
			),
			h.TBody(
				renderRows(p),
			),
		),
		InlineStyle(Assets, "/Table.css"),
	)
}
