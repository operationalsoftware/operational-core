package homeview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"

	g "github.com/maragudk/gomponents"
)

type HomePageProps struct {
	Ctx reqcontext.ReqContext
}

func HomePage(p *HomePageProps) g.Node {

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
