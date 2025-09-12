package galleryview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/reqcontext"
	"slices"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type GalleryPageProps struct {
	Id                int
	Ctx               reqcontext.ReqContext
	Gallery           model.Gallery
	EditURL           string
	AllowedOperations []string
}

func GalleryPage(p *GalleryPageProps) g.Node {

	isEditable := slices.Contains(p.AllowedOperations, "edit")

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
				isEditable,
				components.Button(&components.ButtonProps{
					ButtonType: "primary",
					Classes: c.Classes{
						"edit-button": true,
					},
					Link: p.EditURL,
				},
					components.Icon(&components.IconProps{
						Identifier: "pencil",
					}),
					g.Text("Edit"),
				),
			),
		),
		h.Div(

			// h.Div(
			// 	h.Class("gallery-grid"),
			// g.Group(galleryItems),
			// ),

			components.Gallery(galleryImages),
		),
	})

	return layout.Page(layout.PageProps{
		Title: "Stock Item Details",
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "package-variant-closed",
				Title:          "Gallery",
				URLPart:        "gallery",
			},
			// {Title: stockItem.StockCode},
		},
		Content: content,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/galleryview/gallery_page.css"),
		},
	})
}
