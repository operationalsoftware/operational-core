package router

import (
	"app/internal/handler"
	"app/internal/service"
	"net/http"
)

func addFileRoutes(
	mux *http.ServeMux,
	fileService service.FileService,
) {
	fileHandler := handler.NewFileHandler(fileService)

	mux.HandleFunc("GET /files/{fileID}/complete", fileHandler.CompleteFileUpload)
}
