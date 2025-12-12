package router

import (
	"app/internal/handler"
	"app/internal/service"
	"net/http"
)

func addPDFRoutes(
	mux *http.ServeMux,
	pdfService service.PDFService,
	printNodeService service.PrintNodeService,
) {
	pdfHandler := handler.NewPDFHandler(pdfService, printNodeService)

	mux.HandleFunc("GET /pdf/generate", pdfHandler.PDFGeneratorPage)
	mux.HandleFunc("POST /pdf/generate", pdfHandler.PDFHandler)
	mux.HandleFunc("POST /pdf/print", pdfHandler.PDFPrintHandler)
	mux.HandleFunc("GET /pdf/logs", pdfHandler.PDFGenerationLogsPage)

}
