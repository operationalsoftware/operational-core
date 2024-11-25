package layout

import (
	"app/components"
	"app/internal/db"
	"app/models/usermodel"
	"net/url"
	"strconv"
	"strings"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type breadcrumbsType struct {
	title          string
	renderTitle    func(urlPart string) string
	urlPart        string
	iconIdentifier string
	children       []breadcrumbsType
}

// breadcrumbsDef is a definition of the breadcrumbs
// sorted on urlParts with any wildcards marked as "*"
var breadcrumbsDef = breadcrumbsType{
	title:          "Home",
	urlPart:        "",
	iconIdentifier: "home",

	children: []breadcrumbsType{
		{
			urlPart:        "users",
			title:          "Users",
			iconIdentifier: "account-group",

			children: []breadcrumbsType{
				{
					urlPart:        "add",
					title:          "Add New",
					iconIdentifier: "account-plus",
				},
				{
					urlPart:        "add-api-user",
					title:          "Add API User",
					iconIdentifier: "account-plus",
				},
				{
					urlPart: "*",
					renderTitle: func(urlPart string) string {
						userId, err := strconv.Atoi(urlPart)
						if err != nil {
							return "Error"
						}
						db := db.UseDB()
						user, err := usermodel.ByID(db, userId)
						if err != nil {
							return "Error"
						}
						return user.Username
					},
					iconIdentifier: "account",

					children: []breadcrumbsType{
						{
							urlPart: "edit",
							title:   "Edit",
						},
						{
							urlPart: "reset-password",
							title:   "Reset password",
						},
					},
				},
			},
		},
	},
}

type matchedBreadcrumb struct {
	title          string
	link           string
	iconIdentifier string
}

func breadcrumbs(reqPath string) g.Node {

	// split the url into parts
	pathParts := strings.Split(reqPath, "/")

	// if the url is "/". Split will return an array with two empty parts
	// remove the last one
	if len(pathParts) > 1 && pathParts[len(pathParts)-1] == "" {
		pathParts = pathParts[:len(pathParts)-1]
	}

	// decode/unescape each part
	for i, urlPart := range pathParts {
		unescaped, _ := url.PathUnescape(urlPart)
		pathParts[i] = unescaped
	}

	matchedCrumbs := []matchedBreadcrumb{}
	currentBreadCrumbs := []breadcrumbsType{
		breadcrumbsDef,
	}
	// used to build the link for each crumb based on the urlParts
	currentLink := "/"

	for _, urlPart := range pathParts {
		// find the matching breadCrumb

		var matchingCrumb *breadcrumbsType

		for _, crumb := range currentBreadCrumbs {
			// exact match
			if crumb.urlPart == urlPart {
				matchingCrumb = &crumb
				break
			}
		}

		if matchingCrumb == nil {
			// check for a wildcard match, wildcard match will be last crumb, if any
			if len(currentBreadCrumbs) > 0 &&
				currentBreadCrumbs[len(currentBreadCrumbs)-1].urlPart == "*" {
				matchingCrumb = &currentBreadCrumbs[len(currentBreadCrumbs)-1]
			}
		}

		// no match found - page not found
		if matchingCrumb == nil {
			matchedCrumbs = []matchedBreadcrumb{
				{
					title:          "Home",
					link:           "/",
					iconIdentifier: "home",
				},
				{
					title: "Page not found",
					link:  "/",
				},
			}
			break
		}

		var title string
		if matchingCrumb.title != "" {
			title = matchingCrumb.title
		} else if matchingCrumb.renderTitle != nil {
			title = matchingCrumb.renderTitle(urlPart)
		} else {
			panic("no title or renderTitle function")
		}

		if currentLink != "/" {
			currentLink += "/"
		}
		currentLink += urlPart

		matchedCrumbs = append(matchedCrumbs, matchedBreadcrumb{
			title:          title,
			link:           currentLink,
			iconIdentifier: matchingCrumb.iconIdentifier,
		})

		// recurse
		currentBreadCrumbs = matchingCrumb.children
	}

	index := 0
	crumbNodes := g.Group(
		g.Map(matchedCrumbs, func(c matchedBreadcrumb) g.Node {
			index++

			var crumbContent = g.Group([]g.Node{
				g.If(c.iconIdentifier != "", components.Icon(&components.IconProps{
					Identifier: c.iconIdentifier,
				})),
				h.Span(g.Text(c.title)),
			})

			return h.Li(
				g.If(index > 1, h.Span(h.Class("divider"), g.Text("/"))),
				g.If(index == len(matchedCrumbs), h.Div(crumbContent)),
				g.If(index < len(matchedCrumbs), h.A(
					h.Href(c.link),
					crumbContent,
				)),
			)
		}),
	)

	return h.Nav(
		h.Class("breadcrumbs"),
		h.Aria("label", "breadcrumbs"),
		h.Ol(
			h.Class("breadcrumbs"),
			crumbNodes,
		),
	)
}
