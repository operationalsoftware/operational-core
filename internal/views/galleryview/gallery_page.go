package galleryview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/reqcontext"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type GalleryPageProps struct {
	Id                int
	Ctx               reqcontext.ReqContext
	Gallery           model.Gallery
	EditURL           string
	ParentBreadcrumbs []layout.Breadcrumb
	AllowedOperations []string
}

func GalleryPage(p *GalleryPageProps) g.Node {

	var galleryItems []g.Node
	var galleryImages []string
	for _, gi := range p.Gallery.Items {
		item := h.Div(
			h.Class("gallery-item"),
			h.Img(
				h.Src(gi.DownloadURL),
				h.Alt("gallery image"),
			),
		)

		galleryItems = append(galleryItems, item)
		galleryImages = append(galleryImages, gi.DownloadURL)
	}

	content := g.Group([]g.Node{

		h.Div(
			h.Class("button-container"),

			g.If(
				p.EditURL != "",
				h.A(
					h.Class("button primary edit-button"),
					h.Href(p.EditURL),
					components.Icon(&components.IconProps{
						Identifier: "pencil",
					}),
					g.Text("Edit"),
				),
			),
		),
		components.Gallery(galleryImages),
	})

	breadcrumbs := append([]layout.Breadcrumb{}, p.ParentBreadcrumbs...)
	if len(breadcrumbs) == 0 {
		breadcrumbs = append(breadcrumbs, layout.HomeBreadcrumb)
	}
	breadcrumbs = append(breadcrumbs, layout.Breadcrumb{
		IconIdentifier: "package-variant-closed",
		Title:          "Gallery",
	})

	return layout.Page(layout.PageProps{
		Title:       "Gallery",
		Breadcrumbs: breadcrumbs,
		Content:     content,
		Ctx:         p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/galleryview/gallery_page.css"),
		},
	})
}
