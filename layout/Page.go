package layout

import (
	reqContext "app/reqcontext"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type PageProps struct {
	Title      string
	Content    g.Node
	AppendHead []g.Node
	AppendBody []g.Node
	Ctx        reqContext.ReqContext
}

func Page(p PageProps) g.Node {

	// Construct head
	head := []g.Node{
		h.Meta(h.Charset("utf-8")),
		h.Meta(h.Name("viewport"), h.Content("width=device-width, initial-scale=1")),
		h.Link(h.Rel("manifest"), h.Href("/static/manifest.json")),
		h.Link(h.Rel("stylesheet"), h.Type("text/css"), h.Href("/static/css/reset.css")),
		h.Link(h.Rel("stylesheet"), h.Type("text/css"), h.Href("/static/css/variables.css")),
		h.Link(h.Rel("stylesheet"), h.Type("text/css"), h.Href("/static/css/global.css")),
		h.Link(h.Rel("stylesheet"), h.Type("text/css"), h.Href("/static/css/components.css")),
		h.Link(h.Rel("stylesheet"), h.Type("text/css"), h.Href("/static/css/layout.css")),
		h.Script(h.Type("text/javascript"), h.Src("https://cdn.jsdelivr.net/gh/gnat/surreal/surreal.js")),
		h.Script(h.Type("text/javascript"), h.Src("/static/js/global.js")),
	}

	head = append(head, p.AppendHead...)

	body := []g.Node{
		layout(&layoutProps{
			content: p.Content,
			Ctx:     p.Ctx,
		}),
		h.Script(h.Type("text/javascript"), h.Src("/static/js/htmx.min.js")),
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
