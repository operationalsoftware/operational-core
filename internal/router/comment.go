package router

import (
	"app/internal/service"
	"net/http"
)

func addCommentRoutes(
	mux *http.ServeMux,
	commentService service.CommentService,
) {
	// commentHandler := handler.NewCommentHandler(commentService)

	// mux.HandleFunc("POST /comments/{entity}/{entityId}/add/comment", commentHandler.AddComment)

}
