package router

import (
	"app/internal/handler"
	"app/internal/service"
	"net/http"
)

func addStockItemRoutes(
	mux *http.ServeMux,
	stockItemService service.StockItemService,
) {
	stockItemHandler := handler.NewStockItemHandler(stockItemService)

	mux.HandleFunc("GET /stock-items", stockItemHandler.StockItemsPage)

	mux.HandleFunc("GET /stock-items/add", stockItemHandler.AddStockItemPage)
	mux.HandleFunc("POST /stock-items/add", stockItemHandler.AddStockItem)

	mux.HandleFunc("GET /stock-items/{stockCode}", stockItemHandler.StockItemDetailsPage)

	mux.HandleFunc("GET /stock-items/{stockCode}/edit", stockItemHandler.EditStockItemPage)
	mux.HandleFunc("POST /stock-items/{stockCode}/edit", stockItemHandler.EditStockItem)

	mux.HandleFunc("GET /stock-items/sku-generator", stockItemHandler.GenerateSKU)

	mux.HandleFunc("GET /stock-items/sku-config", stockItemHandler.SKUConfigPage)
	mux.HandleFunc("GET /stock-items/sku-config/add", stockItemHandler.AddSKUItemPage)
	mux.HandleFunc("POST /stock-items/sku-config/add", stockItemHandler.AddSKUItem)
	mux.HandleFunc("DELETE /stock-items/sku-config/{skuField}/{code}", stockItemHandler.DeleteSKUItem)

}
