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

	// mux.HandleFunc("POST /files/{entity}/{entityId}/presign", fileHandler.CreateFile)
	mux.HandleFunc("GET /files/{fileId}/complete", fileHandler.CompleteFileUpload)
}
