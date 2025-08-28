package router

import (
	"app/internal/handler"
	"app/internal/service"
	"net/http"
)

func addStockItemRoutes(
	mux *http.ServeMux,
	stockItemService service.StockItemService,
	commentService service.CommentService,
	fileService service.FileService,
) {
	stockItemHandler := handler.NewStockItemHandler(stockItemService, commentService, fileService)

	mux.HandleFunc("GET /stock-items", stockItemHandler.StockItemsPage)

	mux.HandleFunc("GET /stock-items/add", stockItemHandler.AddStockItemPage)
	mux.HandleFunc("POST /stock-items/add", stockItemHandler.AddStockItem)

	mux.HandleFunc("GET /stock-items/{id}", stockItemHandler.StockItemPage)

	mux.HandleFunc("POST /stock-items/{entityID}/comments", stockItemHandler.AddComment)
	mux.HandleFunc("POST /stock-items/{entityID}/comments/{commentID}/attachment", stockItemHandler.AddAttachment)

	mux.HandleFunc("GET /stock-items/{id}/edit", stockItemHandler.EditStockItemPage)
	mux.HandleFunc("POST /stock-items/{id}/edit", stockItemHandler.EditStockItem)

	mux.HandleFunc("GET /get-stock-codes", stockItemHandler.GetStockCodes)
}
