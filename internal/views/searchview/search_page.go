package searchview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/arrayutil"
	"app/pkg/reqcontext"
	"fmt"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type SearchPageProps struct {
	Ctx              reqcontext.ReqContext
	SearchTerm       string
	SearchEntities   []model.SearchEntity
	SelectedEntities []string
	Results          model.SearchResults
	UserPermissions  model.UserPermissions
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

					searchCheckboxes(p.SearchEntities, p.SelectedEntities),
				),
			),

			h.Div(
				h.Class("search-section"),
				g.Group([]g.Node{
					h.Div(
						h.ID("recent-searches"),

						RecentSearches(p.Results.RecentSearches, p.SearchEntities),
					),
					h.Div(
						h.Class("search-results"),
						SearchResults(p.Results, p.SearchTerm, p.UserPermissions),
					),
				}),
			),
		),
	})

	return layout.Page(layout.PageProps{
		Title:   "Search",
		Content: content,
		Ctx:     p.Ctx,
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "magnify",
				Title:          "Search",
			},
		},
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/searchview/search_page.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/searchview/search_page.js"),
		},
	})
}

func searchCheckboxes(
	searchEntities []model.SearchEntity,
	selectedEntities []string,
) g.Node {
	var entityCheckboxes []g.Node
	for _, entity := range searchEntities {
		c := cases.Title(language.English)
		entityTitle := c.String(entity.Title)

		entityCheckboxes = append(entityCheckboxes,
			h.Label(
				h.Input(
					h.Type("checkbox"),
					h.Value(entity.Name),
					h.Name("E"),
					h.Class("filter-checkbox"),
					g.If(arrayutil.Includes(selectedEntities, entity.Name),
						g.Attr("checked", "checked")),
				),
				g.Text(entityTitle),
			),
		)
	}

	return g.Group(entityCheckboxes)
}

// Recent Searches Section
func RecentSearches(
	terms []model.RecentSearch,
	searchEntities []model.SearchEntity,
) g.Node {
	if len(terms) == 0 {
		return g.Group(nil)
	}

	allowedEntityMap := make(map[string]struct{}, len(searchEntities))
	for _, entity := range searchEntities {
		allowedEntityMap[entity.Name] = struct{}{}
	}

	var items []g.Node
	for _, term := range terms {
		url := fmt.Sprintf("/search?Q=%s", term.SearchTerm)

		for _, entity := range term.SearchEntities {
			if _, ok := allowedEntityMap[entity]; !ok {
				continue
			}
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

func SearchResults(results model.SearchResults, searchTerm string, permissions model.UserPermissions) g.Node {
	var resultSections []g.Node

	if len(results.StockItems) > 0 {
		group := []g.Node{
			h.H3(h.Class("result-type-heading"), g.Text("Stock Item Results")),
			h.Ul(
				h.Class("result-group"),
				g.Group(StockItemResults(results.StockItems, searchTerm)),
			),
		}
		resultSections = append(resultSections, group...)
	}

	if permissions.UserAdmin.Access && len(results.Users) > 0 {
		group := []g.Node{
			h.H3(h.Class("result-type-heading"), g.Text("User Results")),
			h.Ul(
				h.Class("result-group"),
				g.Group(UserResults(results.Users, searchTerm)),
			),
		}
		resultSections = append(resultSections, group...)
	}

	if len(results.Resources) > 0 {
		group := []g.Node{
			h.H3(h.Class("result-type-heading"), g.Text("Resource Results")),
			h.Ul(
				h.Class("result-group"),
				g.Group(ResourceResults(results.Resources, searchTerm)),
			),
		}
		resultSections = append(resultSections, group...)
	}

	if len(results.Services) > 0 {
		group := []g.Node{
			h.H3(h.Class("result-type-heading"), g.Text("Service Results")),
			h.Ul(
				h.Class("result-group"),
				g.Group(ServiceResults(results.Services, searchTerm)),
			),
		}
		resultSections = append(resultSections, group...)
	}

	if len(resultSections) == 0 {
		return h.P(h.Class("placeholder"), g.Text("No Search results."))
	}

	return g.Group(resultSections)
}

func StockItemResults(items []model.StockItemSearchResult, searchTerm string) []g.Node {
	var nodes []g.Node

	for _, item := range items {
		stockItemURL := fmt.Sprintf("/stock-items/%d", item.StockItemID)

		nodes = append(nodes,
			h.A(
				h.Href(stockItemURL),
				h.Li(
					h.Class("search-result-item"),
					h.Div(
						h.Class("title"),
						components.Highlight(item.StockCode, searchTerm),
					),
					h.Div(
						h.Strong(g.Text("Description: ")),
						components.Highlight(item.Description, searchTerm),
					),
				),
			),
		)
	}

	return nodes
}

func UserResults(users []model.UserSearchResult, searchTerm string) []g.Node {
	var nodes []g.Node

	for _, user := range users {
		fullName := strings.TrimSpace(fmt.Sprintf("%v %v", user.FirstName, user.LastName))

		userURL := fmt.Sprintf("/users/%d", user.ID)

		nodes = append(nodes,
			h.A(
				h.Href(userURL),

				h.Li(
					h.Class("search-result-item"),
					h.Div(h.Strong(components.Highlight(fullName, searchTerm))),
					h.Div(
						h.Strong(g.Text("Username: ")),
						components.Highlight(user.Username, searchTerm),
					),
					h.Div(
						h.Strong(g.Text("Email: ")),
						components.Highlight(user.Email, searchTerm),
					),
				),
			),
		)
	}

	return nodes
}

func ResourceResults(resources []model.ResourceSearchResult, searchTerm string) []g.Node {
	var nodes []g.Node

	for _, resource := range resources {
		resourceURL := fmt.Sprintf("/resources/%d", resource.ResourceID)
		serviceTeam := "\u2013"
		if strings.TrimSpace(resource.ServiceOwnershipTeamName) != "" {
			serviceTeam = resource.ServiceOwnershipTeamName
		}

		archivedText := ""
		if resource.IsArchived {
			archivedText = " (Archived)"
		}

		teamNode := g.Node(g.Text(serviceTeam))
		if serviceTeam != "\u2013" {
			teamNode = components.Highlight(serviceTeam, searchTerm)
		}

		nodes = append(nodes,
			h.A(
				h.Href(resourceURL),
				h.Li(
					h.Class("search-result-item"),
					h.Div(
						h.Class("title"),
						components.Highlight(resource.Reference, searchTerm),
						g.Text(archivedText),
					),
					h.Div(
						h.Class("search-result-grid"),
						h.Strong(g.Text("Type: ")),
						h.Span(components.Highlight(resource.Type, searchTerm)),
						h.Strong(g.Text("Service Team: ")),
						h.Span(teamNode),
					),
				),
			),
		)
	}

	return nodes
}

func ServiceResults(services []model.ServiceSearchResult, searchTerm string) []g.Node {
	var nodes []g.Node

	for _, service := range services {
		serviceURL := fmt.Sprintf("/services/%d", service.ResourceServiceID)
		startedAt := service.StartedAt.Format("2006-01-02 15:04:05")
		startedBy := "\u2013"
		if strings.TrimSpace(service.StartedByUsername) != "" {
			startedBy = service.StartedByUsername
		}

		startedByNode := g.Node(g.Text(startedBy))
		if startedBy != "\u2013" {
			startedByNode = components.Highlight(startedBy, searchTerm)
		}

		nodes = append(nodes,
			h.A(
				h.Href(serviceURL),
				h.Li(
					h.Class("search-result-item"),
					h.Div(
						h.Class("title"),
						components.Highlight(service.ResourceReference, searchTerm),
					),
					h.Div(
						h.Class("search-result-grid"),
						h.Strong(g.Text("Resource Type: ")),
						h.Span(components.Highlight(service.ResourceType, searchTerm)),
						h.Strong(g.Text("Status: ")),
						h.Span(components.Highlight(string(service.Status), searchTerm)),
						h.Strong(g.Text("Started By: ")),
						h.Span(startedByNode),
						h.Strong(g.Text("Started At: ")),
						h.Span(g.Text(startedAt)),
					),
				),
			),
		)
	}

	return nodes
}
