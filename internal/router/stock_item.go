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
	galleryService service.GalleryService,
	appHMAC apphmac.AppHMAC,
) {
	stockItemHandler := handler.NewStockItemHandler(stockItemService, commentService, galleryService, appHMAC)

	mux.HandleFunc("GET /stock-items", stockItemHandler.StockItemsPage)

	mux.HandleFunc("GET /stock-items/add", stockItemHandler.AddStockItemPage)
	mux.HandleFunc("POST /stock-items/add", stockItemHandler.AddStockItem)

	mux.HandleFunc("GET /stock-items/{id}", stockItemHandler.StockItemPage)

	mux.HandleFunc("GET /stock-items/{id}/edit", stockItemHandler.EditStockItemPage)
	mux.HandleFunc("POST /stock-items/{id}/edit", stockItemHandler.EditStockItem)

	mux.HandleFunc("GET /get-stock-codes", stockItemHandler.GetStockCodes)
}
