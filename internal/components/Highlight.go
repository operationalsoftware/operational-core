package components

import (
	"strings"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

func Highlight(text, term string) g.Node {
	// Add docs for this func
	if term == "" {
		return g.Text(text)
	}

	// Case-insensitive match
	termLower := strings.ToLower(term)
	textLower := strings.ToLower(text)

	start := strings.Index(textLower, termLower)
	if start == -1 {
		return g.Text(text)
	}

	end := start + len(term)

	return g.Group([]g.Node{
		g.Text(text[:start]),
		h.Span(h.Class("highlight"), g.Text(text[start:end])),
		g.Text(text[end:]),
	})
}
