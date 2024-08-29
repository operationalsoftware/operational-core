package home

import (
	"app/components"
	"app/internal/reqcontext"
	"app/layout"

	g "github.com/maragudk/gomponents"
)

type homePageProps struct {
	Ctx reqcontext.ReqContext
}

func homePage(p *homePageProps) g.Node {

	content := components.Card(
		components.GridMenu(&components.GridMenuProps{
			Items:           layout.TopLevelMenuItems,
			UserPermissions: p.Ctx.User.Permissions,
		}),
	)

	return layout.Page(layout.PageProps{
		Title:   "Home",
		Content: content,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/routes/home/home.css"),
		},
	})
}
