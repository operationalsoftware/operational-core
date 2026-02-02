package handler

import (
	"app/internal/model"
	"app/internal/pdftemplate"
	"app/internal/service"
	"app/internal/views/pdfview"
	"app/pkg/printnode"
	"app/pkg/reqcontext"
	"fmt"
	"log"
	"net/http"
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

func (h *PDFHandler) PDFHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/pdf/generate", http.StatusSeeOther)
}

func (h *PDFHandler) PDFGeneratorPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	templates := pdftemplate.SortedTemplates()

	templateName := r.URL.Query().Get("TemplateName")

	var selectedTemplate *pdftemplate.RegisteredTemplate
	if templateName != "" {
		t, found := pdftemplate.Registry[templateName]
		if found {
			selectedTemplate = &t
		}
	}

	pdfview.PDFGeneratorPage(pdfview.PDFPageProps{
		Ctx:              ctx,
		Templates:        templates,
		SelectedTemplate: selectedTemplate,
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

	pdfBuf, resolvedTitle, err := h.pdfService.GenerateFromJSON(
		r.Context(),
		templateName,
		[]byte(inputData),
	)
	if err != nil {
		log.Println("An error occurred generating PDF:", err)
		http.Error(w, "PDF generation failed", http.StatusInternalServerError)
		return
	}

	downloadName := h.pdfService.GeneratePDFFilename(resolvedTitle)

	_, err = h.pdfService.RecordGeneration(
		r.Context(),
		templateName,
		inputData,
		pdfBuf,
		ctx.User.UserID,
		resolvedTitle,
	)
	if err != nil {
		log.Println("An error occurred recording PDF generation log:", err)
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", downloadName))
	w.Write(pdfBuf)
}

func (h *PDFHandler) PDFPrintHandler(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.Printing.Operator {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	printLogIDStr := r.FormValue("PrintLogID")
	requirementName := r.FormValue("RequirementName")
	templateName := r.FormValue("TemplateName")
	inputData := r.FormValue("InputData")
	overridePrinterName := r.FormValue("PrinterName")

	var err error
	if printLogIDStr != "" {
		printLogID, convErr := strconv.Atoi(printLogIDStr)
		if convErr != nil {
			http.Error(w, "Invalid PrintLogID", http.StatusBadRequest)
			return
		}
		_, err = h.pdfService.Reprint(r.Context(), printLogID, overridePrinterName, ctx.User.UserID)
	} else {
		if templateName == "" || inputData == "" || strings.TrimSpace(requirementName) == "" {
			http.Error(w, "TemplateName, InputData, and RequirementName are required", http.StatusBadRequest)
			return
		}
		_, _, err = h.pdfService.PrintAndLog(r.Context(), templateName, inputData, requirementName, ctx.User.UserID)
	}

	if err != nil {
		log.Println("An error occurred printing PDF:", err)
		http.Error(w, "Print failed", http.StatusInternalServerError)
		return
	}

	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
		return
	}
	ref := r.Referer()
	if ref == "" {
		ref = "/printing"
	}
	http.Redirect(w, r, ref, http.StatusSeeOther)
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

// PrinterAssignmentsPage shows the mapping of print requirements to printers.
func (h *PDFHandler) PrinterAssignmentsPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.Printing.Admin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	assignments, err := h.pdfService.ListPrintRequirements(r.Context())
	if err != nil {
		log.Println("An error occurred fetching print requirements:", err)
	}

	pdfview.PrinterAssignmentsPage(pdfview.PrinterAssignmentsPageProps{
		Ctx:         ctx,
		Assignments: assignments,
	}).Render(w)
}

// PDFPrintJobsPage shows print jobs.
func (h *PDFHandler) PDFPrintJobsPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.Printing.Operator {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	printNodeStatus, err := h.printNodeService.Status(r.Context())
	if err != nil {
		log.Println("An error occurred checking PrintNode status:", err)
	}
	printers := []printnode.Printer{}
	if printNodeStatus.Configured {
		if list, err := h.printNodeService.Printers(r.Context()); err == nil {
			printers = list
		}
	}

	printLogs, err := h.pdfService.ListRecentPrintLogs(r.Context(), 10)
	if err != nil {
		log.Println("An error occurred fetching PDF print logs:", err)
	}

	pdfview.PDFPrintJobsPage(pdfview.PDFPrintJobsPageProps{
		Ctx:             ctx,
		PrintLogs:       printLogs,
		Printers:        printers,
		PrintNodeStatus: printNodeStatus,
	}).Render(w)
}

// PrinterAssignmentEditPage shows an edit form for a print requirement.
func (h *PDFHandler) PrinterAssignmentEditPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.Printing.Admin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	reqName := strings.TrimSpace(r.URL.Query().Get("RequirementName"))
	if reqName == "" {
		http.Error(w, "RequirementName is required", http.StatusBadRequest)
		return
	}

	ctx = reqcontext.GetContext(r)
	printNodeStatus, err := h.printNodeService.Status(r.Context())
	if err != nil {
		log.Println("An error occurred checking PrintNode status:", err)
		http.Error(w, "Unable to check PrintNode status", http.StatusInternalServerError)
		return
	}
	printers := []printnode.Printer{}
	if printNodeStatus.Configured {
		if list, err := h.printNodeService.Printers(r.Context()); err == nil {
			printers = list
		}
	}

	assignments, err := h.pdfService.ListPrintRequirements(r.Context())
	if err != nil {
		log.Println("An error occurred fetching print requirements:", err)
	}
	availablePrinters, _ := h.pdfService.ListAvailablePrinters(r.Context(), reqName, printers)
	var current model.PrintRequirement
	for _, a := range assignments {
		if strings.EqualFold(a.RequirementName, reqName) {
			current = a
			break
		}
	}

	pdfview.PrinterAssignmentEditPage(pdfview.PrinterAssignmentEditPageProps{
		Ctx:            ctx,
		Requirement:    reqName,
		Printers:       availablePrinters,
		SelectedName:   current.PrinterName,
		PrintNodeReady: printNodeStatus.Configured,
	}).Render(w)
}

// PrinterAssignmentsSave updates a requirement → printer mapping.
func (h *PDFHandler) PrinterAssignmentsSave(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.Printing.Admin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	requirement := r.FormValue("RequirementName")
	printerName := strings.TrimSpace(r.FormValue("PrinterName"))

	if requirement == "" {
		http.Error(w, "RequirementName is required", http.StatusBadRequest)
		return
	}
	if printerName == "" {
		http.Error(w, "PrinterName is required", http.StatusBadRequest)
		return
	}

	_, err := h.pdfService.SavePrintRequirement(r.Context(), requirement, printerName, ctx.User.UserID)
	if err != nil {
		log.Println("An error occurred saving print requirement:", err)
		err := fmt.Errorf("Failed to assign printer %w", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
		return
	}
	http.Redirect(w, r, "/printing/printer-assignments", http.StatusSeeOther)
}
