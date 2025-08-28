package handler

import (
	"app/internal/model"
	"app/internal/service"
	"app/pkg/validate"
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

	return ve
}
