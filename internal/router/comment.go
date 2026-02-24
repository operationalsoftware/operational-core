package router

import (
	"app/internal/handler"
	"app/internal/service"
	"app/pkg/apphmac"
	"net/http"
)

func addCommentRoutes(
	mux *http.ServeMux,
	commentService service.CommentService,
	userService service.UserService,
	fileService service.FileService,
	appHMAC apphmac.AppHMAC,
) {
	commentHandler := handler.NewCommentHandler(commentService, userService, fileService, appHMAC)

	// Centralized comment endpoints using comment_thread_id
	mux.HandleFunc("POST /comments/{threadID}/add", commentHandler.Add)
	mux.HandleFunc("POST /comments/{threadID}/{commentID}/attachment", commentHandler.AddAttachment)
	mux.HandleFunc("GET /comments/mentions/users", commentHandler.MentionUsers)
}
