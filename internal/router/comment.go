package router

import (
	"app/internal/handler"
	"app/internal/service"
	"net/http"
)

func addCommentRoutes(
	mux *http.ServeMux,
	commentService service.CommentService,
	fileService service.FileService,
	hmacService service.HMACService,
) {
	commentHandler := handler.NewCommentHandler(commentService, fileService, hmacService)

	// Centralized comment endpoints using comment_thread_id
	mux.HandleFunc("POST /comments/{threadID}/add", commentHandler.Add)
	mux.HandleFunc("POST /comments/{threadID}/{commentID}/attachment", commentHandler.AddAttachment)
}
