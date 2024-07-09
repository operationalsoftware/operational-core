package homeviews

import (
	"app/components"
	"app/layout"
	"app/reqcontext"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type HomePageProps struct {
	Ctx reqcontext.ReqContext
}

func HomePage(p *HomePageProps) g.Node {

	content := g.Group([]g.Node{
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
	})

	return layout.Page(layout.PageProps{
		Title:   "Home",
		Content: content,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/src/indexView.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/src/index.js"),
		},
	})
}
