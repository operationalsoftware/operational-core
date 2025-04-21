package searchview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"
	"strings"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type SearchPageProps struct {
	Ctx reqcontext.ReqContext
}

func SliceToSet(slice []string) map[string]bool {
	set := make(map[string]bool, len(slice))
	for _, v := range slice {
		set[v] = true
	}
	return set
}

func SearchPage(p SearchPageProps) g.Node {
	searchTerm := p.Ctx.Req.URL.Query().Get("q")
	searchTypes := strings.Split(p.Ctx.Req.URL.Query().Get("types"), ",")
	typesMap := SliceToSet(searchTypes)

	content := g.Group([]g.Node{
		h.Div(
			h.Class("search-wrapper"),

			components.Form(
				h.ID("search-form"),

				h.Input(
					h.Type("text"),
					h.Name("q"),
					h.ID("search-input"),
					h.Placeholder("Search"),
					h.Value(searchTerm),
					h.AutoComplete("off"),
				),
			),

			h.Div(
				h.Class("filters"),
				h.Label(
					h.Input(
						h.Type("checkbox"),
						h.Value("batch"),
						h.Class("filter-checkbox"),
						g.If(
							typesMap["batch"],
							g.Attr("checked", "checked"),
						),
					),
					g.Text("Batch"),
				),
				h.Label(
					h.Input(
						h.Type("checkbox"),
						h.Value("user"),
						h.Class("filter-checkbox"),
						g.If(
							typesMap["user"],
							g.Attr("checked", "checked"),
						),
					),
					g.Text("User"),
				),
			),

			h.Div(
				h.Class("search-results"),
			),
		),
	})

	return layout.Page(layout.PageProps{
		Title:   "Search",
		Content: content,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/searchview/search_page.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/searchview/search_page.js"),
		},
	})
}
