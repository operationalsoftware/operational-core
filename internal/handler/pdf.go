package handler

import (
	"app/internal/pdftemplate"
	"app/internal/service"
	"app/internal/views/pdfview"
	"app/pkg/printnode"
	"app/pkg/reqcontext"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type PDFHandler struct {
	pdfService       service.PDFService
	printNodeService service.PrintNodeService
}

func NewPDFHandler(
	pdfService service.PDFService,
	printNodeService service.PrintNodeService,
) *PDFHandler {
	return &PDFHandler{
		pdfService:       pdfService,
		printNodeService: printNodeService,
	}
}

func (h *PDFHandler) PDFGeneratorPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	printNodeStatus, err := h.printNodeService.Status(r.Context())
	if err != nil {
		log.Println("An error occurred checking PrintNode status:", err)
	}

	templates := pdftemplate.SortedTemplates()

	var printers []printnode.Printer
	if printNodeStatus.Configured {
		printers, err = h.printNodeService.Printers(r.Context())
		if err != nil {
			log.Println("An error occurred fetching PrintNode printers:", err)
		}
	}

	defaultPrinterID := ""
	defaultPrinterName := ""
	for _, p := range printers {
		if p.Default {
			defaultPrinterID = fmt.Sprintf("%d", p.ID)
			defaultPrinterName = p.Name
			break
		}
	}
	if defaultPrinterID == "" && len(printers) > 0 {
		defaultPrinterID = fmt.Sprintf("%d", printers[0].ID)
		defaultPrinterName = printers[0].Name
	}

	printRequirements := make([]pdfview.PrintRequirement, 0, len(templates))
	for _, tmpl := range templates {
		printRequirements = append(printRequirements, pdfview.PrintRequirement{
			Name:               tmpl.Name,
			Description:        tmpl.Description,
			SelectedPrinter:    defaultPrinterID,
			DefaultPrinterName: defaultPrinterName,
		})
	}

	templateName := r.URL.Query().Get("TemplateName")

	var selectedTemplate *pdftemplate.RegisteredTemplate
	if templateName != "" {
		t, found := pdftemplate.Registry[templateName]
		if found {
			selectedTemplate = &t
		}
	}

	logs, err := h.pdfService.ListRecentLogs(r.Context(), 10)
	if err != nil {
		log.Println("An error occurred fetching PDF generation logs:", err)
	}

	printLogs, err := h.pdfService.ListRecentPrintLogs(r.Context(), 10)
	if err != nil {
		log.Println("An error occurred fetching PDF print logs:", err)
	}

	pdfview.PDFGeneratorPage(pdfview.PDFPageProps{
		Ctx:               ctx,
		Templates:         templates,
		PrintNodeStatus:   printNodeStatus,
		Printers:          printers,
		PrintRequirements: printRequirements,
		SelectedTemplate:  selectedTemplate,
		GenerationLogs:    logs,
		PrintLogs:         printLogs,
	}).Render(w)
}

func (h *PDFHandler) PDFHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		log.Println("An error occurred parsing form:", err)
		http.Error(w, "PDF generation failed", http.StatusInternalServerError)
		return
	}

	templateName := r.FormValue("TemplateName")
	inputData := r.FormValue("InputData")

	ctx := reqcontext.GetContext(r)

	pdfBuf, err := h.pdfService.GenerateFromJSON(r.Context(),
		templateName,
		[]byte(inputData))
	if err != nil {
		log.Println("An error occurred generating PDF:", err)
		http.Error(w, "PDF generation failed", http.StatusInternalServerError)
		return
	}

	_, err = h.pdfService.RecordGeneration(
		r.Context(),
		templateName,
		inputData,
		pdfBuf,
		ctx.User.UserID,
	)
	if err != nil {
		log.Println("An error occurred recording PDF generation log:", err)
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline")
	w.Write(pdfBuf)
}

func (h *PDFHandler) PDFPrintHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx := reqcontext.GetContext(r)

	printLogIDStr := r.FormValue("PrintLogID")
	printerIDStr := r.FormValue("PrinterID")
	printerName := r.FormValue("PrinterName")
	requirementName := r.FormValue("RequirementName")
	templateName := r.FormValue("TemplateName")
	inputData := r.FormValue("InputData")

	printerID, _ := strconv.Atoi(printerIDStr)

	var err error
	if printLogIDStr != "" {
		printLogID, convErr := strconv.Atoi(printLogIDStr)
		if convErr != nil {
			http.Error(w, "Invalid PrintLogID", http.StatusBadRequest)
			return
		}
		_, err = h.pdfService.Reprint(r.Context(), printLogID, printerID, printerName, ctx.User.UserID)
	} else {
		if templateName == "" || inputData == "" || printerID == 0 {
			http.Error(w, "TemplateName, InputData, and PrinterID are required", http.StatusBadRequest)
			return
		}
		_, err = h.pdfService.PrintAndLog(r.Context(), templateName, inputData, printerID, printerName, requirementName, ctx.User.UserID)
	}

	if err != nil {
		log.Println("An error occurred printing PDF:", err)
		http.Error(w, "Print failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}
