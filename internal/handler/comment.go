package handler

import (
	"app/internal/model"
	"app/internal/service"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type CommentHandler struct {
	commentService service.CommentService
}

func NewCommentHandler(
	commentService service.CommentService,
) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

func (h *CommentHandler) AddComment(w http.ResponseWriter, r *http.Request) {
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

	var fd addAndonCommentFormData
	if err := json.NewDecoder(r.Body).Decode(&fd); err != nil {
		log.Println(err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	fd.normalise()

	entityID, err := strconv.Atoi(entityIDStr)
	if err != nil {
		http.Error(w, "Invalid entity Id", http.StatusBadRequest)
		return
	}

	commentId, err := h.commentService.CreateComment(
		r.Context(),
		&model.NewComment{
			Comment:  fd.Comment,
			Entity:   entity,
			EntityID: entityID,
		},
		ctx.User.UserID,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error creating andon", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"commentId": commentId,
	})

	// http.Redirect(w, r, fmt.Sprintf("/andons/%d", fd.EntityID), http.StatusSeeOther)
}

type addAndonCommentFormData struct {
	Comment     string
	EntityID    string
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

	// if fd.EntityID == 0 {
	// 	ve.Add("EntityID", "is required")
	// }

	return ve
}
