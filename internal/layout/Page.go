package layout

import (
	"app/pkg/reqcontext"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
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
		h.Link(h.Rel("stylesheet"), h.Type("text/css"), h.Href("/static/css/reset.css")),
		h.Link(h.Rel("stylesheet"), h.Type("text/css"), h.Href("/static/css/variables.css")),
		h.Link(h.Rel("stylesheet"), h.Type("text/css"), h.Href("/static/css/global.css")),
		h.Link(h.Rel("stylesheet"), h.Type("text/css"), h.Href("/static/css/components.css")),
		h.Link(h.Rel("stylesheet"), h.Type("text/css"), h.Href("/static/css/layout.css")),
		h.Script(h.Type("text/javascript"), h.Src("/static/js/global.js")),
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
