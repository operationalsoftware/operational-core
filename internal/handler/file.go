package handler

import (
	"app/internal/service"
	"app/pkg/reqcontext"
	"log"
	"net/http"
)

type FileHandler struct {
	fileService service.FileService
}

func NewFileHandler(fileService service.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

func (h *FileHandler) CompleteFileUpload(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	if !ctx.User.Permissions.UserAdmin.Access {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	fileID := r.PathValue("fileID")

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

	w.WriteHeader(http.StatusOK)
}
