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
	icons := []layout.Item{
		{
			Icon: "account-group",
			Name: "Users",
			Link: "/users",
			Role: []string{"User Admin", "User"},
		},
		{
			Icon: "github",
			Name: "Github",
			Link: "/github",
			Role: []string{"User", "Dummy"},
		},
	}

	indexContent := g.Group([]g.Node{
		h.H1(g.Text("Operational Core Home")),
		components.InlineScript("/src/index.js"),

		components.Card(
			h.Div(
				h.Class("modules-container"),
				h.Div(
					h.Class("module-items"),
					g.Group(g.Map(icons, func(i layout.Item) g.Node {
						for _, role := range i.Role {
							if utils.CheckRole(p.Ctx.User.Roles, role) {
								return h.A(
									h.Class("module-item"),
									h.Href(i.Link),
									components.Icon(&components.IconProps{
										Identifier: i.Icon,
									}),
									h.Div(
										h.Class("module-item-name"),
										g.Text(i.Name),
									),
								)
							}
						}
						return g.Text("")
					})),
				),
			),
		),
		components.InlineStyle("/src/indexView.css"),
	})

	return layout.Page(layout.PageProps{
		Title:   "Home",
		Content: indexContent,
		Ctx:     p.Ctx,
	})
}
