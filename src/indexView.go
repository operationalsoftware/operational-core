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
		h.FormEl(
			h.Method("GET"),
			components.Select(&components.SelectProps{
				Options: []components.Option{
					{Label: "Option 1", Value: "1"},
					{Label: "Option 2", Value: "2"},
					{Label: "Option 3", Value: "3"},
				},
				Value:    []string{"2"},
				Name:     "select",
				ID:       "select",
				Multiple: false,
			}),
			components.Button(&components.ButtonProps{
				ButtonType: components.ButtonPrimary,
			},
				h.Type("submit"),
				g.Text("Submit"),
			),
		),
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
