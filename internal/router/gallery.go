package router

import (
	"app/internal/handler"
	"app/internal/service"
	"app/pkg/apphmac"
	"net/http"
)

func addGalleryRoutes(
	mux *http.ServeMux,
	fileService service.FileService,
	galleryService service.GalleryService,
	appHMAC apphmac.AppHMAC,
) {
	galleryHandler := handler.NewGalleryHandler(fileService, galleryService, appHMAC)

	mux.HandleFunc("GET /gallery/{galleryID}", galleryHandler.GalleryPage)
	mux.HandleFunc("GET /gallery/{galleryID}/edit", galleryHandler.EditPage)
	mux.HandleFunc("POST /gallery/{galleryID}/edit", galleryHandler.AddGalleryItem)
	mux.HandleFunc("PUT /gallery/{galleryID}/{galleryItemID}/set-position", galleryHandler.SetGalleryItemPosition)
	mux.HandleFunc("DELETE /gallery/{galleryID}/item/{galleryItemID}", galleryHandler.DeleteGalleryItem)
}
