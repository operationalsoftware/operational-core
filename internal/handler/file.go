package handler

import (
	"app/internal/model"
	"app/internal/pdftemplate"
	"app/internal/service"
	"app/internal/views/pdfview"
	"app/pkg/reqcontext"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type FileHandler struct {
	fileService service.FileService
}

func NewFileHandler(fileService service.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

func (h *FileHandler) PDFGeneratorPage(w http.ResponseWriter, r *http.Request) {
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

func (h *FileHandler) CreateFile(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	entity := r.PathValue("entity")
	entityIDStr := r.PathValue("entityId")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		http.Error(w, "Invalid entity Id", http.StatusBadRequest)
		return
	}

	var fd addFileFormData
	if err := json.NewDecoder(r.Body).Decode(&fd); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	file, signedURL, err := h.fileService.CreateFile(
		r.Context(),
		&model.File{
			Filename:    fd.Filename,
			ContentType: fd.ContentType,
			SizeBytes:   fd.SizeBytes,
			Entity:      entity,
			EntityID:    entityID,
		},
		ctx.User.UserID,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error adding file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"fileId":    file.FileID,
		"signedUrl": signedURL,
	})

}

type addFileFormData struct {
	Filename    string
	ContentType string
	SizeBytes   int
	Entity      string
	EntityID    int
}

func (h *FileHandler) CompleteFileUpload(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	fileID := r.PathValue("fileId")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	err := h.fileService.CompleteFileUpload(
		r.Context(), fileID,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error updating file status", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"fileId": fileID,
		"status": "success",
	})

}
