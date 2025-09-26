package handler

import (
	"app/internal/pdftemplate"
	"app/internal/service"
	"app/internal/views/pdfview"
	"app/pkg/reqcontext"
	"log"
	"net/http"
)

type PDFHandler struct {
	pdfService service.PDFService
}

func NewPDFHandler(pdfService service.PDFService) *PDFHandler {
	return &PDFHandler{pdfService: pdfService}
}

func (h *PDFHandler) PDFGeneratorPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	props := pdfview.PDFPageProps{
		Ctx:       ctx,
		Templates: pdftemplate.SortedTemplates(),
	}

	templateName := r.URL.Query().Get("TemplateName")

	if templateName != "" {
		t, found := pdftemplate.Registry[templateName]
		if found {
			props.SelectedTemplate = &t
		}
	}

	pdfview.PDFGeneratorPage(props).Render(w)
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

	pdfBuf, err := h.pdfService.GenerateFromJSON(r.Context(), templateName, []byte(inputData))
	if err != nil {
		log.Println("An error occurred generating PDF:", err)
		http.Error(w, "PDF generation failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline")
	w.Write(pdfBuf)
}
