package searchview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/arrayutil"
	"app/pkg/reqcontext"
	"fmt"
	"strings"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type SearchPageProps struct {
	Ctx     reqcontext.ReqContext
	Results map[string][]model.SearchResult
}

func SearchPage(p SearchPageProps) g.Node {
	searchTerm := p.Ctx.Req.URL.Query().Get("Q")
	searchEntities := p.Ctx.Req.URL.Query()["E"]

	content := g.Group([]g.Node{
		h.Div(
			h.Class("search-wrapper"),

			components.Form(
				h.Method("GET"),
				h.ID("search-form"),

				h.Div(
					h.Class("search-input"),

					h.Input(
						h.Type("text"),
						h.Name("Q"),
						h.ID("search-input"),
						h.Placeholder("Search"),
						h.Value(searchTerm),
						h.AutoComplete("off"),
					),

					h.Button(
						h.Type("submit"),
						h.Class("button"),
						g.Text("Search"),
					),
				),

				h.Div(
					h.Class("filters"),
					h.Label(
						h.Input(
							h.Type("checkbox"),
							h.Value("batch"),
							h.Name("E"),
							h.Class("filter-checkbox"),
							g.Attr("data-type", "batch"),
							g.If(
								arrayutil.Includes(searchEntities, "batch"),
								g.Attr("checked", "checked"),
							),
						),
						g.Text("Batch"),
					),
					h.Label(
						h.Input(
							h.Type("checkbox"),
							h.Value("user"),
							h.Name("E"),
							h.Class("filter-checkbox"),
							g.Attr("data-type", "user"),
							g.If(
								arrayutil.Includes(searchEntities, "user"),
								g.Attr("checked", "checked"),
							),
						),
						g.Text("User"),
					),
				),
			),

			h.Div(
				h.Class("search-results"),
				renderSearchResults(p.Results),
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

func renderResultItems(resultType string, items []model.SearchResult) []g.Node {
	var nodes []g.Node

	for _, item := range items {
		switch resultType {
		case "user":
			dataMap, ok := item.Data.(model.UserSearchResult)
			if !ok {
				continue // skip if not in expected format
			}

			fullName := strings.TrimSpace(fmt.Sprintf("%v %v", dataMap.FirstName, dataMap.LastName))

			nodes = append(nodes,
				h.Li(
					h.Class("search-result-item"),
					h.Div(h.Strong(g.Text(fullName))),
					h.Div(
						h.Strong(g.Text("Username: ")),
						g.Text(fmt.Sprintf("%v", dataMap.Username)),
					),
					h.Div(
						h.Strong(g.Text("Email: ")),
						g.Text(fmt.Sprintf("%v", dataMap.Email)),
					),
				),
			)

		case "batch":
			dataMap, ok := item.Data.(model.BatchSearchResult)
			if !ok {
				continue
			}

			nodes = append(nodes,
				h.Li(
					h.Class("search-result-item"),
					h.Div(g.Text("Batch #: "+fmt.Sprintf("%v", dataMap.BatchNumber))),
					h.Div(g.Text("Works Order #: "+fmt.Sprintf("%v", dataMap.WorksOrderNumber))),
					h.Div(g.Text("Part #: "+fmt.Sprintf("%v", dataMap.PartNumber))),
				),
			)
		}
	}

	return nodes
}

func renderSearchResults(results map[string][]model.SearchResult) g.Node {
	if len(results) == 0 {
		return h.P(h.Class("placeholder"), g.Text("No Search results."))
	}

	var resultSections []g.Node

	for resultType, items := range results {
		if len(items) == 0 {
			continue
		}

		title := strings.Title(resultType)

		group := []g.Node{
			h.H3(h.Class("result-type-heading"), g.Text(title+" Results")),
			h.Ul(
				h.Class("result-group"),
				g.Group(renderResultItems(resultType, items)),
			),
		}

		resultSections = append(resultSections, group...)
	}

	if len(resultSections) == 0 {
		return h.P(h.Class("placeholder"), g.Text("No Search results."))
	}

	return g.Group(resultSections)
}
