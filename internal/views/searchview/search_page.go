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
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type SearchPageProps struct {
	Ctx            reqcontext.ReqContext
	SearchTerm     string
	SearchEntities []string
	Results        model.SearchResults
}

func SearchPage(p SearchPageProps) g.Node {

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
						h.Value(p.SearchTerm),
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

					SearchCheckboxes(p.SearchEntities),
				),
			),

			h.Div(
				h.Class("search-section"),
				g.Group([]g.Node{
					h.Div(
						h.ID("recent-searches"),

						RecentSearches(p.Results.RecentSearches),
					),
					h.Div(
						h.Class("search-results"),
						SearchResults(p.Results),
					),
				}),
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

func SearchCheckboxes(searchEntities []string) g.Node {
	searchEntitiesList := []string{"user"}

	var entityCheckboxes []g.Node
	for _, entity := range searchEntitiesList {
		c := cases.Title(language.English)
		entityTitle := c.String(entity)

		fmt.Println(entity)
		fmt.Println(entityTitle)

		entityCheckboxes = append(entityCheckboxes,
			h.Label(
				h.Input(
					h.Type("checkbox"),
					h.Value(entity),
					h.Name("E"),
					h.Class("filter-checkbox"),
					g.If(arrayutil.Includes(searchEntities, entity),
						g.Attr("checked", "checked")),
				),
				g.Text(entityTitle),
			),
		)
	}

	return g.Group(entityCheckboxes)
}

// Recent Searches Section
func RecentSearches(terms []model.RecentSearch) g.Node {
	if len(terms) == 0 {
		return g.Group(nil)
	}

	var items []g.Node
	for _, term := range terms {
		url := fmt.Sprintf("/search?Q=%s", term.SearchTerm)

		for _, entity := range term.SearchEntities {
			url += fmt.Sprintf("&E=%s", entity)
		}

		items = append(items,
			h.Li(
				h.A(
					h.Href(url),
					g.Text(term.SearchTerm),
				),
			),
		)

	}

	return h.Div(
		h.Class("recent-searches"),
		h.H3(h.Class("title"), g.Text("Recent Searches")),
		h.Ul(
			h.Class("list"),
			g.Group(items),
		),
	)
}

func SearchResults(results model.SearchResults) g.Node {
	var resultSections []g.Node

	if len(results.Users) > 0 {
		group := []g.Node{
			h.H3(h.Class("result-type-heading"), g.Text("User Results")),
			h.Ul(
				h.Class("result-group"),
				g.Group(UserResults(results.Users)),
			),
		}
		resultSections = append(resultSections, group...)
	}

	if len(resultSections) == 0 {
		return h.P(h.Class("placeholder"), g.Text("No Search results."))
	}

	return g.Group(resultSections)
}

func UserResults(users []model.UserSearchResult) []g.Node {
	var nodes []g.Node

	for _, user := range users {
		fullName := strings.TrimSpace(fmt.Sprintf("%v %v", user.FirstName, user.LastName))

		nodes = append(nodes,
			h.Li(
				h.Class("search-result-item"),
				h.Div(h.Strong(g.Text(fullName))),
				h.Div(
					h.Strong(g.Text("Username: ")),
					g.Text(user.Username),
				),
				h.Div(
					h.Strong(g.Text("Email: ")),
					g.Text(user.Email),
				),
			),
		)
	}

	return nodes
}
