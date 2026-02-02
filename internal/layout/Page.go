package layout

import (
	"app/pkg/reqcontext"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type PageProps struct {
	Title             string
	Breadcrumbs       []Breadcrumb
	Content           g.Node
	LayoutMainPadding *bool
	AppendHead        []g.Node
	AppendBody        []g.Node
	Ctx               reqcontext.ReqContext
}

func Page(p PageProps) g.Node {

	// Construct head
	head := []g.Node{
		h.Link(h.Rel("manifest"), h.Href("/static/manifest.json")),
		h.Link(h.Rel("icon"), h.Href("/static/favicon/32x32.png"), h.Type("image/png")),
		h.Link(h.Rel("apple-touch-icon"), h.Href("/static/favicon/apple-touch-icon.png"), g.Attr("sizes", "180x180")),
		h.Link(h.Rel("stylesheet"), h.Type("text/css"), h.Href("/static/css/reset.css")),
		h.Link(h.Rel("stylesheet"), h.Type("text/css"), h.Href("/static/css/colours.css")),
		h.Link(h.Rel("stylesheet"), h.Type("text/css"), h.Href("/static/css/variables.css")),
		h.Link(h.Rel("stylesheet"), h.Type("text/css"), h.Href("/static/css/global.css")),
		h.Link(h.Rel("stylesheet"), h.Type("text/css"), h.Href("/static/css/components.css")),
		h.Link(h.Rel("stylesheet"), h.Type("text/css"), h.Href("/static/css/layout.css")),
		h.Script(h.Type("text/javascript"), h.Src("/static/js/global.js")),

		h.Meta(h.Name("viewport"), h.Content("width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no")),
		// iOS Meta Tags for Web App
		h.Meta(h.Name("apple-mobile-web-app-title"), h.Content("Operational Core")),

		// Use mobile-web-app-capable for modern browsers (cross-platform)
		h.Meta(h.Name("mobile-web-app-capable"), h.Content("yes")),

		// Keep the iOS-specific meta tag for legacy support (still needed for iOS)
		h.Meta(h.Name("apple-mobile-web-app-capable"), h.Content("yes")),
		h.Meta(h.Name("apple-mobile-web-app-status-bar-style"), h.Content("black-translucent")),
	}

	head = append(head, p.AppendHead...)

	body := []g.Node{
		h.Div(h.ID("loading-message"), g.Text("Loading...")),
		layout(&layoutProps{
			breadcrumbs: p.Breadcrumbs,
			content:     p.Content,
			mainPadding: p.LayoutMainPadding,
			ctx:         p.Ctx,
		}),
	}

	body = append(body, p.AppendBody...)

	// HTML5 boilerplate document
	return c.HTML5(c.HTML5Props{
		Title:       p.Title,
		Description: "",
		Language:    "en",
		Head:        head,
		Body:        body,
	})
}
