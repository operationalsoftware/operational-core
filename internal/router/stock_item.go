package router

import (
	"app/internal/handler"
	"app/internal/service"
	"app/pkg/apphmac"
	"net/http"
)

func addStockItemRoutes(
	mux *http.ServeMux,
	stockItemService service.StockItemService,
	commentService service.CommentService,
	fileService service.FileService,
	galleryService service.GalleryService,
	appHMAC apphmac.AppHMAC,
) {
	stockItemHandler := handler.NewStockItemHandler(stockItemService, commentService, fileService, galleryService, appHMAC)

	mux.HandleFunc("GET /stock-items", stockItemHandler.StockItemsPage)

	mux.HandleFunc("GET /stock-items/add", stockItemHandler.AddStockItemPage)
	mux.HandleFunc("POST /stock-items/add", stockItemHandler.AddStockItem)

	mux.HandleFunc("GET /stock-items/{id}", stockItemHandler.StockItemPage)

	mux.HandleFunc("GET /stock-items/{id}/edit", stockItemHandler.EditStockItemPage)
	mux.HandleFunc("POST /stock-items/{id}/edit", stockItemHandler.EditStockItem)

	mux.HandleFunc("GET /get-stock-codes", stockItemHandler.GetStockCodes)
}
