package homeview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"

	g "maragu.dev/gomponents"
)

type HomePageProps struct {
	Ctx reqcontext.ReqContext
}

func HomePage(p *HomePageProps) g.Node {

	content := components.Card(
		components.GridMenu(&components.GridMenuProps{
			Groups:          layout.AppMenu,
			UserPermissions: p.Ctx.User.Permissions,
		}),
	)

	return layout.Page(layout.PageProps{
		Ctx:         p.Ctx,
		Title:       "Home",
		Breadcrumbs: []layout.Breadcrumb{layout.HomeBreadcrumb},
		Content:     content,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/homeview/home_page.css"),
		},
	})
}
