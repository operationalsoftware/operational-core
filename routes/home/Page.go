package home

import (
	"app/components"
	"app/internal/reqcontext"
	"app/layout"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type homePageProps struct {
	Ctx reqcontext.ReqContext
}

func homePage(p *homePageProps) g.Node {

	content := g.Group([]g.Node{
		components.Card(
			h.Div(
				h.Class("modules-container"),
				h.Div(
					h.Class("module-items"),
					g.Group(g.Map(layout.ModuleGridItems, func(m layout.ModuleGridItem) g.Node {

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
			components.InlineStyle("/routes/home/home.css"),
		},
	})
}
