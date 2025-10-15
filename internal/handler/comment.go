package handler

import (
	"app/internal/model"
	"app/internal/service"
	"app/pkg/apphmac"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CommentHandler struct {
	commentService service.CommentService
	fileService    service.FileService
	appHMAC        apphmac.AppHMAC
}

func NewCommentHandler(
	commentService service.CommentService,
	fileService service.FileService,
	appHMAC apphmac.AppHMAC,
) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		fileService:    fileService,
		appHMAC:        appHMAC,
	}
}

type addCommentFormData struct {
	Comment string
}

func (fd *addCommentFormData) normalise() {
	fd.Comment = strings.TrimSpace(fd.Comment)
}

func (fd *addCommentFormData) validate() validate.ValidationErrors {
	var ve validate.ValidationErrors = make(map[string][]string)

	if fd.Comment == "" {
		ve.Add("Comment", "should not be empty")
	}

	return ve
}

func (h *CommentHandler) Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	threadIDStr := r.PathValue("threadID")
	threadID, err := strconv.Atoi(threadIDStr)
	if err != nil || threadID <= 0 {
		http.Error(w, "Invalid thread id", http.StatusBadRequest)
		return
	}

	// Decode request body once into a struct that can carry the envelope
	var reqBody struct {
		Comment string `json:"comment"`
		HMAC    string `json:"hmac"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Verify HMAC: prefer envelope
	isValid, err := h.appHMAC.CheckEnvelope(reqBody.HMAC, "comment_thread", fmt.Sprintf("%d", threadID), "add")
	if err != nil || !isValid {
		http.Error(w, "Error validating", http.StatusUnauthorized)
		return
	}

	// Validate comment content
	fd := addCommentFormData{Comment: reqBody.Comment}
	fd.normalise()
	if ve := fd.validate(); len(ve) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{"errors": ve})
		return
	}

	// Create the comment
	ctx := reqcontext.GetContext(r)
	commentID, err := h.commentService.CreateComment(r.Context(), &model.NewComment{
		CommentThreadID: threadID,
		Comment:         fd.Comment,
	}, ctx.User.UserID)
	if err != nil {
		http.Error(w, "Error creating comment", http.StatusInternalServerError)
		return
	}

	// Also return an attachment envelope so the client can upload files for this new comment
	attachPayload := apphmac.Payload{
		Entity:      "comment",
		EntityID:    fmt.Sprintf("%d", commentID),
		Permissions: []string{"add"},
		Expires:     time.Now().Add(24 * time.Hour).Unix(),
	}
	attachEnvelope := h.appHMAC.CreateEnvelope(
		attachPayload,
	)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"commentId":      commentID,
		"attachmentHmac": attachEnvelope,
	})
}

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
		Filename    string `json:"filename"`
		ContentType string `json:"contentType"`
		SizeBytes   int    `json:"sizeBytes"`
		HMAC        string `json:"hmac"`
	}
	if err := json.NewDecoder(r.Body).Decode(&fd); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	isValid, err := h.appHMAC.CheckEnvelope(fd.HMAC, "comment", fmt.Sprintf("%d", commentID), "add")
	if err != nil || !isValid {
		http.Error(w, "Error validating", http.StatusUnauthorized)
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
