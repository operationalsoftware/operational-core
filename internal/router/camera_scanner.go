package router

import (
	"app/internal/views/camerascannerview"
	"app/pkg/reqcontext"
	"net/http"
)

func addCameraScannerRoutes(
	mux *http.ServeMux,
) {

	mux.HandleFunc("GET /camera-scanner", func(w http.ResponseWriter, r *http.Request) {
		ctx := reqcontext.GetContext(r)

		_ = camerascannerview.CameraScannerApp(&camerascannerview.CameraScannerAppProps{
			Ctx: ctx,
		}).Render(w)
	})

}
