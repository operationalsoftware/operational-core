package router

import (
	"app/internal/handler"
	"app/internal/service"
	"net/http"
)

func addGalleryRoutes(
	mux *http.ServeMux,
	fileService service.FileService,
	galleryService service.GalleryService,
) {
	galleryHandler := handler.NewGalleryHandler(fileService, galleryService)

	mux.HandleFunc("GET /gallery/{galleryID}", galleryHandler.GalleryPage)
	mux.HandleFunc("GET /gallery/{galleryID}/edit", galleryHandler.EditPage)
	mux.HandleFunc("POST /gallery/{galleryID}/edit", galleryHandler.AddGalleryItem)
	mux.HandleFunc("PUT /gallery/{galleryID}/reorder", galleryHandler.ReorderGalleryItem)
	mux.HandleFunc("DELETE /gallery/{galleryID}/item/{galleryItemID}", galleryHandler.DeleteGalleryItem)
}
