package handler

import (
	"app/internal/pdftemplate"
	"app/internal/service"
	"app/internal/views/pdfview"
	"app/pkg/printnode"
	"app/pkg/reqcontext"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

	pdfTitle := h.pdfService.GeneratePDFTitleFromInput(templateName, []byte(inputData))
	downloadName := h.pdfService.GeneratePDFFilename(pdfTitle)

	pdfBuf, err := h.pdfService.GenerateFromJSON(
		r.Context(),
		templateName,
		[]byte(inputData),
		pdfTitle,
	)
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
		pdfTitle,
	)
	if err != nil {
		log.Println("An error occurred recording PDF generation log:", err)
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", downloadName))
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

// PDFViewHandler renders a simple inline viewer page for a given PDF URL.
func (h *PDFHandler) PDFViewHandler(w http.ResponseWriter, r *http.Request) {
	resolvedURL, resolvedTitle, err := h.resolvePDFView(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pdfview.PDFInlineViewer(pdfview.PDFInlineViewerProps{
		Title: resolvedTitle,
		Src:   resolvedURL,
	}).Render(w)
}

// PDFStreamHandler streams a PDF inline with an inline content disposition to avoid forced downloads.
func (h *PDFHandler) PDFStreamHandler(w http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Query().Get("file_id")
	if fileID == "" {
		http.Error(w, "missing file_id", http.StatusBadRequest)
		return
	}

	data, filename, err := h.pdfService.FetchPDFFile(r.Context(), fileID)
	if err != nil {
		log.Println("An error occurred streaming PDF:", err)
		http.Error(w, "PDF not found", http.StatusNotFound)
		return
	}

	if filename == "" {
		filename = "document.pdf"
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", template.HTMLEscapeString(filename)))
	w.Write(data)
}

func (h *PDFHandler) resolvePDFView(r *http.Request) (string, string, error) {
	rawURL := r.URL.Query().Get("url")
	fileID := r.URL.Query().Get("file_id")
	title := r.URL.Query().Get("title")

	if fileID != "" {
		rawURL = fmt.Sprintf("/pdf/stream?file_id=%s", url.QueryEscape(fileID))
		if fetchedTitle, err := h.pdfService.GetPDFTitleByFileID(r.Context(), fileID); err == nil && strings.TrimSpace(fetchedTitle) != "" {
			title = fetchedTitle
		} else if err != nil {
			log.Println("An error occurred fetching PDF title:", err)
		}
	}

	if strings.TrimSpace(rawURL) == "" {
		return "", "", fmt.Errorf("missing url")
	}
	if strings.TrimSpace(title) == "" {
		title = "PDF Document"
	}

	return rawURL, title, nil
}

// PDFGenerationLogsPage renders a paginated table of PDF generation logs.
func (h *PDFHandler) PDFGenerationLogsPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	page, _ := strconv.Atoi(r.URL.Query().Get("Page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("PageSize"))
	if pageSize <= 0 {
		pageSize = 25
	}
	offset := (page - 1) * pageSize

	logs, total, err := h.pdfService.ListGenerationLogs(r.Context(), pageSize, offset)
	if err != nil {
		log.Println("An error occurred fetching PDF generation logs:", err)
		http.Error(w, "Unable to load logs", http.StatusInternalServerError)
		return
	}

	pdfview.PDFGenerationLogsPage(pdfview.PDFGenerationLogsPageProps{
		Ctx:       ctx,
		Logs:      logs,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		PageQuery: "Page",
		SizeQuery: "PageSize",
	}).Render(w)
}
