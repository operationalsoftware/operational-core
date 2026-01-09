package router

import (
	"app/internal/handler"
	"app/internal/service"
	"net/http"
)

func addPrintingRoutes(
	mux *http.ServeMux,
	pdfService service.PDFService,
	printNodeService service.PrintNodeService,
) {
	printingHandler := handler.NewPDFHandler(pdfService, printNodeService)

	mux.HandleFunc("GET /printing", printingHandler.PDFPrintJobsPage)
	mux.HandleFunc("GET /printing/print-jobs", printingHandler.PDFPrintJobsPage)
	mux.HandleFunc("GET /printing/printer-assignments", printingHandler.PrinterAssignmentsPage)
	mux.HandleFunc("GET /printing/printer-assignments/edit", printingHandler.PrinterAssignmentEditPage)
	mux.HandleFunc("POST /printing/printer-assignments", printingHandler.PrinterAssignmentsSave)
}
