package handler

import (
	"app/internal/service"
	"app/internal/views/pdfview"
	"app/pkg/reqcontext"
	"encoding/json"
	"net/http"
)

type FileHandler struct {
	pdfService service.PDFService
}

func NewFileHandler(pdfService service.PDFService) *FileHandler {
	return &FileHandler{pdfService: pdfService}
}

func (h *FileHandler) PDFGeneratorPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	pdfview.PDFGeneratorPage(pdfview.PDFPageProps{
		Ctx:       ctx,
		Templates: []string{"invoice", "receipt"},
	}).Render(w)

	return
}

func (h *PDFHandler) FileHandler(w http.ResponseWriter, r *http.Request) {

	templateName := r.FormValue("template")
	rawJSON := r.FormValue("params")

	var inputParams map[string]interface{}
	if err := json.Unmarshal([]byte(rawJSON), &inputParams); err != nil {
		http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
		return
	}

	pdfBuf, err := h.pdfService.GeneratePDF(r.Context(), templateName, inputParams)
	if err != nil {
		http.Error(w, "PDF generation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline")
	w.Write(pdfBuf)
}
