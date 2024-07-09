package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

func LogoImgFull(classes c.Classes, children ...g.Node) g.Node {
	if classes == nil {
		classes = c.Classes{}
	}

	classes["logo-img-full"] = true

	// copy classes
	darkClasses := c.Classes{}
	for key, value := range classes {
		darkClasses[key] = value
	}

	darkClasses["dark-theme"] = true

	return g.Group([]g.Node{
		h.Img(
			classes,
			h.Src("/static/img/rms-logo.png"),
			g.Group(children),
		),
		h.Img(
			darkClasses,
			h.Src("/static/img/rms-logo-dark-theme.png"),
			g.Group(children),
		),
	})
}
