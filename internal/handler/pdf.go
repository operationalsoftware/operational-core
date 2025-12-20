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

func (h *PDFHandler) PDFGeneratorPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	printNodeStatus, err := h.printNodeService.Status(r.Context())
	if err != nil {
		log.Println("An error occurred checking PrintNode status:", err)
	}

	templates := pdftemplate.SortedTemplates()

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

	pdfview.PDFGeneratorPage(pdfview.PDFPageProps{
		Ctx:              ctx,
		Templates:        templates,
		PrintNodeStatus:  printNodeStatus,
		SelectedTemplate: selectedTemplate,
		GenerationLogs:   logs,
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
		"",
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
	if printerID != 0 && strings.TrimSpace(printerName) == "" {
		if printers, err := h.printNodeService.Printers(r.Context()); err == nil {
			for _, pr := range printers {
				if pr.ID == printerID {
					printerName = pr.Name
					break
				}
			}
		}
	}

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

	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
		return
	}
	ref := r.Referer()
	if ref == "" {
		ref = "/pdf/printer-assignments"
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

// PrinterAssignmentsPage shows the mapping of print requirements to printers plus recent print logs.
func (h *PDFHandler) PrinterAssignmentsPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

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

	assignments, err := h.pdfService.ListPrintRequirements(r.Context())
	if err != nil {
		log.Println("An error occurred fetching print requirements:", err)
	}
	printLogs, err := h.pdfService.ListRecentPrintLogs(r.Context(), 10)
	if err != nil {
		log.Println("An error occurred fetching PDF print logs:", err)
	}

	pdfview.PrinterAssignmentsPage(pdfview.PrinterAssignmentsPageProps{
		Ctx:         ctx,
		Printers:    printers,
		Assignments: assignments,
		PrintLogs:   printLogs,
	}).Render(w)
}

// PrinterAssignmentEditPage shows an edit form for a print requirement.
func (h *PDFHandler) PrinterAssignmentEditPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.Automation.PrinterAssignmentsEditor {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	reqName := strings.TrimSpace(r.URL.Query().Get("RequirementName"))
	if reqName == "" {
		http.Error(w, "RequirementName is required", http.StatusBadRequest)
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
		SelectedID:     current.PrinterID,
		PrintNodeReady: printNodeStatus.Configured,
	}).Render(w)
}

// PrinterAssignmentsSave updates a requirement → printer mapping.
func (h *PDFHandler) PrinterAssignmentsSave(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.Automation.PrinterAssignmentsEditor {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	requirement := r.FormValue("RequirementName")
	printerIDStr := r.FormValue("PrinterID")

	if requirement == "" {
		http.Error(w, "RequirementName is required", http.StatusBadRequest)
		return
	}
	printerID, _ := strconv.Atoi(printerIDStr)

	printerName := ""
	if printerID > 0 {
		if printers, err := h.printNodeService.Printers(r.Context()); err == nil {
			for _, pr := range printers {
				if pr.ID == printerID {
					printerName = pr.Name
					break
				}
			}
		}
	}

	_, err := h.pdfService.SavePrintRequirement(r.Context(), model.PrintRequirement{
		RequirementName: requirement,
		PrinterID:       printerID,
		PrinterName:     printerName,
		AssignedBy:      ctx.User.UserID,
	})
	if err != nil {
		log.Println("An error occurred saving print requirement:", err)
		if strings.Contains(err.Error(), "printer already assigned") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "Save failed", http.StatusInternalServerError)
		return
	}

	if strings.Contains(r.Header.Get("Accept"), "application/json") {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
		return
	}
	http.Redirect(w, r, "/pdf/printer-assignments", http.StatusSeeOther)
}
