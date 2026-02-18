package router

import (
	"app/internal/handler"
	"net/http"
)

func addImageToTextRoutes(
	mux *http.ServeMux,
) {
	handler := handler.NewImageToTextHandler()
	mux.HandleFunc("GET /image-to-text/resolve", handler.ImageToTextResolvePage)

}
