package src

import (
	"app/components"
	"app/layout"
	"app/utils"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type indexViewProps struct {
	Ctx utils.Context
}

func indexView(p *indexViewProps) g.Node {

	indexContent := g.Group([]g.Node{
		components.Card(
			h.Div(
				h.Class("modules-container"),
				h.Div(
					h.Class("module-items"),
					g.Group(g.Map(layout.AppGalleryModules, func(m layout.AppGalleryModule) g.Node {

						if m.Show != nil {
							show := m.Show(p.Ctx.User.Permissions)
							if !show {
								return g.Text("")
							}
						}

						return h.A(
							h.Class("module-item"),
							h.Href(m.Link),
							components.Icon(&components.IconProps{
								Identifier: m.Icon,
							}),
							h.Div(
								h.Class("module-item-name"),
								g.Text(m.Name),
							),
						)
					})),
				),
			),
		),
		components.InlineScript("/src/index.js"),
		components.InlineStyle("/src/indexView.css"),
	})

	return layout.Page(layout.PageProps{
		Title:   "Home",
		Content: indexContent,
		Ctx:     p.Ctx,
	})
}
