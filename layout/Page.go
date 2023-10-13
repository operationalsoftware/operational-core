package layout

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

func script(src string) g.Node {
	return Script(Type("text/javascript"), Src(src))
}

func stylesheet(src string) g.Node {
	return Link(Rel("stylesheet"), Type("text/css"), Href(src))
}

type PageParams struct {
	Title   string
	Crumbs  []Crumb
	Content g.Node
	Scripts []string
	CSS     []string
}

func Page(params PageParams) g.Node {

	// Construct head
	head := []g.Node{
		Meta(Charset("utf-8")),
		Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
		Link(Rel("manifest"), Href("/manifest.json")),
	}

	// Add common css
	cssUrls := []string{
		"/css/reset.css",
		"/css/variables.css",
		"/css/global.css",
	}
	head = append(head, g.Map(cssUrls, stylesheet)...)

	// Add additional css
	if params.CSS != nil {
		head = append(head, g.Map(params.CSS, stylesheet)...)
	}

	// Add common scripts
	scriptUrls := []string{
		"https://cdn.jsdelivr.net/gh/gnat/surreal/surreal.js",
		"https://cdn.jsdelivr.net/gh/gnat/css-scope-inline/script.js",
	}
	head = append(head, g.Map(scriptUrls, script)...)

	// Add additional scripts
	if params.Scripts != nil {
		head = append(head, g.Map(params.Scripts, script)...)
	}

	// HTML5 boilerplate document
	return c.HTML5(c.HTML5Props{
		Title:       params.Title,
		Description: "",
		Language:    "en",
		Head:        head,
		Body: layout(
			layoutParams{
				content: params.Content,
				crumbs:  params.Crumbs,
			},
		),
	})
}
