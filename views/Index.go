package views

import (
	"fmt"
	o "operationalcore/components"
	"operationalcore/layout"
	"operationalcore/utils"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type CustomDataRow struct {
	col1 string
	col2 string
	col3 string
	col4 int
}

func (t CustomDataRow) Render() map[string]o.RenderedCell {
	return map[string]o.RenderedCell{
		"col-1": {
			Content: g.Text(t.col1),
			Classes: c.Classes{
				"table-cell": true,
			},
		},
		"col-2": {
			Content: g.Text(t.col2),
			Classes: c.Classes{
				"table-cell": true,
			},
		},
		"col-3": {
			Content: g.Text(t.col3),
			Classes: c.Classes{
				"table-cell": true,
			},
		},
		"col-4": {
			Content: g.Text(fmt.Sprint(t.col4)),
			Classes: c.Classes{
				"table-cell": true,
			},
		},
	}
}

var columns = []o.TableColumn{
	{
		Name: "Column 1",
		Key:  "col-1",
	},
	{
		Name: "Column 2",
		Key:  "col-2",
	},
	{
		Name: "Column 3",
		Key:  "col-3",
	},
	{
		Name: "Column 4",
		Key:  "col-4",
	},
}

var props = &o.TableProps{
	Columns: columns,
	Data:    data,
}

var data = []o.TableRowRenderer{
	CustomDataRow{
		col1: "Data 1",
		col2: "Data 2",
		col3: "Data 3",
		col4: 4,
	},
	CustomDataRow{
		col1: "Data 1",
		col2: "Data 2",
		col3: "Data 3",
		col4: 8,
	},
}

type IndexProps struct {
	Ctx utils.Context
}

func createIndexCrumbs() []layout.Crumb {
	var crumbs = []layout.Crumb{
		{
			Title:    "Home",
			LinkPart: "",
			Icon:     "home",
		},
	}

	return crumbs
}

func Index(p *IndexProps) g.Node {

	crumbs := createIndexCrumbs()

	indexContent := g.Group([]g.Node{
		h.H1(g.Text("Operational Core Home")),
		o.InlineScript(Assets, "/Index.js"),
	})

	return layout.Page(layout.PageProps{
		Title:   "Home",
		Content: indexContent,
		Ctx:     p.Ctx,
		Crumbs:  crumbs,
	})
}
