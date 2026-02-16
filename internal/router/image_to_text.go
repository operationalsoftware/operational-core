package router

import (
	"app/internal/views/imagetotextview"
	"app/pkg/reqcontext"
	"net/http"
)

func addImageToTextRoutes(
	mux *http.ServeMux,
) {

	mux.HandleFunc("GET /image-to-text", func(w http.ResponseWriter, r *http.Request) {
		ctx := reqcontext.GetContext(r)

		_ = imagetotextview.ImageToTextPage(&imagetotextview.ImageToTextPageProps{
			Ctx: ctx,
		}).Render(w)
	})

	mux.HandleFunc("GET /image-to-text/fix", func(w http.ResponseWriter, r *http.Request) {
		ctx := reqcontext.GetContext(r)

		_ = imagetotextview.ImageToTextFixPage(&imagetotextview.ImageToTextFixPageProps{
			Ctx: ctx,
		}).Render(w)
	})

}
