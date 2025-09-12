package galleryview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/reqcontext"
	"fmt"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

type EditGalleryPageProps struct {
	Id                int
	Ctx               reqcontext.ReqContext
	Gallery           model.Gallery
	GalleryID         int
	HMAC              string
	AllowedOperations []string
	Expires           int64
	GalleryPageURL    string
}

func EditGalleryPage(p *EditGalleryPageProps) g.Node {

	var galleryItemsNodes []g.Node
	if len(p.Gallery.Items) > 0 {
		for _, item := range p.Gallery.Items {
			node := h.Div(
				g.Attr("class", "gallery-item"),
				g.Attr("data-id", fmt.Sprintf("%d", item.GalleryItemID)),
				g.Attr("data-gallery-id", fmt.Sprintf("%d", item.GalleryID)),

				h.Div(
					h.Class("drag-handle"),
					components.Icon(&components.IconProps{
						Classes: c.Classes{
							"icon": true,
						},
						Identifier: "drag-horizontal-variant",
					}),

					h.Button(
						g.Attr("class", "btn-danger"),
						g.Attr("onclick", fmt.Sprintf("deleteItem(%d, %d, %d)", p.Gallery.GalleryID, item.GalleryItemID, item.Position)),

						components.Icon(&components.IconProps{
							Identifier: "close",
						}),
					),
				),

				h.Div(
					h.Class("thumbnail"),
					h.Style(fmt.Sprintf("background-image: url(%s)", item.DownloadURL)),

					// h.Img(h.Src(item.DownloadURL), g.Attr("class", "thumbnail")),
				),
			)
			galleryItemsNodes = append(galleryItemsNodes, node)
		}
	}

	content := g.Group([]g.Node{

		h.Div(

			h.Div(
				h.Class("header"),

				h.Div(
					h.Class("title"),
					h.H3(g.Text("Edit Gallery")),
					h.P(g.Text("Add, reorder, or remove items in this gallery.")),
				),

				h.Form(
					h.Class("gallery-edit-form"),
					h.Name("gallery-edit-form"),
					h.Method("POST"),
					h.EncType("multipart/form-data"),
					g.Attr("onsubmit", "submitGalleryItems(event)"),

					h.Div(
						h.Class("comment-box"),

						h.Input(
							h.Name("EntityID"),
							h.Type("hidden"),
							h.Value(fmt.Sprintf("%d", p.GalleryID)),
						),

						h.Div(
							h.Class("submit-btn"),
							components.Button(&components.ButtonProps{
								Classes: c.Classes{
									"upload-gallery-btn": true,
								},
								ButtonType: "primary",
								Loading:    true,
								Disabled:   true,
							},
								h.Type("submit"),
								components.Icon(&components.IconProps{
									Identifier: "upload",
								}),
								g.Text(" Upload Items"),
							),
						),

						h.Div(
							h.Class("file-upload-wrapper"),

							h.Div(
								h.Class("files"),

								h.Div(
									h.ID("selected-files"),
									h.Class("selected-files"),
								),

								h.Label(
									h.Class("file-input-label button small"),

									h.Input(
										h.Class("file-input"),
										h.Name("Files"),
										h.Type("file"),
										h.Multiple(),
										h.Accept("image/*,video/*,application/pdf"),
										g.Attr("data-max-files", "10"),
									),

									components.Icon(&components.IconProps{
										Identifier: "paperclip-plus",
									}),
									g.Text("Choose files"),
								),
							),
						),
					),
				),
			),

			g.If(len(p.Gallery.Items) == 0,
				h.Div(
					g.Attr("class", "empty-state"),
					h.P(g.Text("No items in this gallery yet. Start by uploading files.")),
				),
			),
			g.If(len(p.Gallery.Items) > 0,
				h.Div(
					g.Attr("class", "gallery-grid"),
					g.Attr("data-gallery-id", fmt.Sprintf("%d", p.GalleryID)),
					g.Group(galleryItemsNodes),
				),
			),
		),
	})

	return layout.Page(layout.PageProps{
		Title: "Edit Gallery",
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "package-variant-closed",
				Title:          "Gallery",
				URLPart:        fmt.Sprintf("gallery/%d?%s", p.GalleryID, p.GalleryPageURL),
			},
			{
				IconIdentifier: "pencil",
				Title:          "Edit",
			},
		},
		Content: content,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/galleryview/edit_page.css"),
			components.InlineScript("/internal/views/galleryview/edit_page.js"),
		},
	})
}
