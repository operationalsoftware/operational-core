package router

import (
	"app/internal/handler"
	"app/internal/service"
	"net/http"
)

func addPDFRoutes(
	mux *http.ServeMux,
	pdfService service.PDFService,
) {
	pdfHandler := handler.NewPDFHandler(pdfService)

	mux.HandleFunc("GET /pdf/generate", pdfHandler.PDFGeneratorPage)
	mux.HandleFunc("POST /pdf/generate", pdfHandler.PDFHandler)

}
