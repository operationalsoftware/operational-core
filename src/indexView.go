package src

import (
	"app/components"
	"app/layout"
	reqContext "app/reqcontext"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type indexViewProps struct {
	Ctx reqContext.ReqContext
}

func indexView(p *indexViewProps) g.Node {

	indexContent := g.Group([]g.Node{
		components.SearchSelect(&components.SearchSelectProps{
			Options: []components.Option{
				{
					Label: "Option 1",
					Value: "option1",
				},
				{
					Label: "Option 2",
					Value: "option2",
				},
			},
			Multiple: true,
		}),

		components.SearchSelect(&components.SearchSelectProps{
			Options: []components.Option{
				{
					Label: "Option 1",
					Value: "option1",
				},
				{
					Label: "Option 2",
					Value: "option2",
				},
			},
			Multiple: false,
		}),
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
		Content: indexContent,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/src/indexView.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/src/index.js"),
		},
	})
}
