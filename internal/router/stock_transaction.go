package router

import (
	"app/internal/handler"
	"app/internal/service"
	"net/http"
)

func addStockTrxRoutes(
	mux *http.ServeMux,
	stockTrxService service.StockTrxService,
) {
	stockTrxHandler := handler.NewStockTrxHandler(stockTrxService)

	// Stock Home page
	mux.HandleFunc("GET /stock", stockTrxHandler.StockLevelsPage)

	// Stock transactions page
	mux.HandleFunc("GET /stock/transactions", stockTrxHandler.StockTransactionsPage)

	// Post transaction pages
	mux.HandleFunc("GET /stock/post-transaction", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/stock/post-transaction/stock-movement", http.StatusSeeOther)
	})

	// Stock Movement
	mux.HandleFunc("GET /stock/post-transaction/stock-movement", stockTrxHandler.PostStockMovementPage)
	mux.HandleFunc("POST /stock/post-transaction/stock-movement", stockTrxHandler.PostStockMovement)

	// Production
	mux.HandleFunc("GET /stock/post-transaction/production", stockTrxHandler.PostProductionPage)
	mux.HandleFunc("POST /stock/post-transaction/production", stockTrxHandler.PostProduction)

	// Production Reversal
	mux.HandleFunc("GET /stock/post-transaction/production-reversal", stockTrxHandler.PostProductionReversalPage)
	mux.HandleFunc("POST /stock/post-transaction/production-reversal", stockTrxHandler.PostProductionReversal)

	// Consumption
	mux.HandleFunc("GET /stock/post-transaction/consumption", stockTrxHandler.PostConsumptionPage)
	mux.HandleFunc("POST /stock/post-transaction/consumption", stockTrxHandler.PostConsumption)

	// Consumption Reversal
	mux.HandleFunc("GET /stock/post-transaction/consumption-reversal", stockTrxHandler.PostConsumptionReversalPage)
	mux.HandleFunc("POST /stock/post-transaction/consumption-reversal", stockTrxHandler.PostConsumptionReversal)

}
