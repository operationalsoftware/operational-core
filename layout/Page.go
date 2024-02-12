package layout

import (
	"operationalcore/utils"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

func script(src string) g.Node {
	return h.Script(h.Type("text/javascript"), h.Src(src))
}

func stylesheet(src string) g.Node {
	return h.Link(h.Rel("stylesheet"), h.Type("text/css"), h.Href(src))
}

type PageProps struct {
	Title   string
	Content g.Node
	Scripts []string
	CSS     []string
	Ctx     utils.Context
}

func Page(p PageProps) g.Node {

	// Construct head
	head := []g.Node{
		h.Meta(h.Charset("utf-8")),
		h.Meta(h.Name("viewport"), h.Content("width=device-width, initial-scale=1")),
		h.Link(h.Rel("manifest"), h.Href("/static/manifest.json")),
	}

	// Add common css
	cssUrls := []string{
		"/static/css/reset.css",
		"/static/css/variables.css",
		"/static/css/global.css",
	}
	head = append(head, g.Map(cssUrls, stylesheet)...)

	// Add additional css
	if p.CSS != nil {
		head = append(head, g.Map(p.CSS, stylesheet)...)
	}

	// Add common scripts
	scriptUrls := []string{
		"https://cdn.jsdelivr.net/gh/gnat/surreal/surreal.js",
		"https://cdn.jsdelivr.net/gh/gnat/css-scope-inline/script.js",
		"/static/js/htmx.min.js",
		"/static/js/global.js",
	}
	head = append(head, g.Map(scriptUrls, script)...)

	// Add additional scripts
	if p.Scripts != nil {
		head = append(head, g.Map(p.Scripts, script)...)
	}

	// HTML5 boilerplate document
	return c.HTML5(c.HTML5Props{
		Title:       p.Title,
		Description: "",
		Language:    "en",
		Head:        head,
		Body: layout(
			&layoutProps{
				content: p.Content,
				Ctx:     p.Ctx,
			},
		),
	})
}
