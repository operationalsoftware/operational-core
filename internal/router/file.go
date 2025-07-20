package router

import (
	"app/internal/handler"
	"app/internal/service"
	"net/http"
)

func addFileRoutes(
	mux *http.ServeMux,
	pdfService service.PDFService,
) {
	pdfHandler := handler.NewPDFHandler(pdfService)

	mux.HandleFunc("PUT /file/generate", pdfHandler.PDFGeneratorPage)
	// mux.HandleFunc("POST /pdf/generate", pdfHandler.PDFHandler)

}
