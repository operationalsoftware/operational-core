package handler

import (
	"app/internal/model"
	"app/internal/service"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type CommentHandler struct {
	commentService service.CommentService
	fileService    service.FileService
}

func NewCommentHandler(
	commentService service.CommentService,
	fileService service.FileService,
) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		fileService:    fileService,
	}
}

type addAndonCommentFormData struct {
	Comment     string
	Attachments []model.File
}

func (fd *addAndonCommentFormData) normalise() {
	fd.Comment = strings.TrimSpace(fd.Comment)
}

func (fd *addAndonCommentFormData) validate() validate.ValidationErrors {
	var ve validate.ValidationErrors = make(map[string][]string)

	if fd.Comment == "" {
		ve.Add("Comment", "should not be empty")
	}

	return ve
}

// Add handles POST /comments/{threadID}/add
func (h *CommentHandler) Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	threadIDStr := r.PathValue("threadID")
	threadID, err := strconv.Atoi(threadIDStr)
	if err != nil {
		http.Error(w, "Invalid thread id", http.StatusBadRequest)
		return
	}

	var fd addAndonCommentFormData
	if err := json.NewDecoder(r.Body).Decode(&fd); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	fd.normalise()
	if ve := fd.validate(); len(ve) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{"errors": ve})
		return
	}

	ctx := reqcontext.GetContext(r)
	commentID, err := h.commentService.CreateComment(r.Context(), &model.NewComment{
		CommentThreadID: threadID,
		Comment:         fd.Comment,
	}, ctx.User.UserID)
	if err != nil {
		http.Error(w, "Error creating comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"commentId": commentID})
}

// AddAttachment handles POST /comments/{commentID}/attachment
func (h *CommentHandler) AddAttachment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	commentIDStr := r.PathValue("commentID")
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment id", http.StatusBadRequest)
		return
	}

	var fd struct {
		Filename    string
		ContentType string
		SizeBytes   int
	}
	if err := json.NewDecoder(r.Body).Decode(&fd); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	ctx := reqcontext.GetContext(r)
	file, signedURL, err := h.fileService.CreateFile(
		r.Context(),
		&model.File{
			Filename:    fd.Filename,
			ContentType: fd.ContentType,
			SizeBytes:   fd.SizeBytes,
			Entity:      "Comment",
			EntityID:    commentID,
		},
		ctx.User.UserID,
	)
	if err != nil {
		http.Error(w, "Error adding file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"fileId":    file.FileID,
		"signedUrl": signedURL,
	})
}
