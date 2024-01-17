package components

import (
	"io"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type IconProps struct {
	Identifier string
	Classes    c.Classes
}

func Icon(p *IconProps) g.Node {
	if p.Classes == nil {
		p.Classes = c.Classes{}
	}

	p.Classes["icon"] = true

	file, err := Assets.Open("/icon-svgs/" + p.Identifier + ".svg")
	if err != nil {
		return g.Text("")
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	// Convert the byte slice to a string
	svgString := string(fileBytes)

	// svg component
	svg := h.SVG(
		p.Classes,
		g.Raw(svgString),
	)

	return svg
}
