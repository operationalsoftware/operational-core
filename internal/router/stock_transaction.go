package router

import (
	"app/internal/handler"
	"app/internal/service"
	"net/http"
)

func addStockTransactionRoutes(
	mux *http.ServeMux,
	stockTransactionService service.StockTransactionService,
) {
	stockTransactionHandler := handler.NewStockTransactionHandler(stockTransactionService)

	// Stock Home page
	mux.HandleFunc("GET /stock", stockTransactionHandler.StockLevelsPage)

	// Stock transactions page
	mux.HandleFunc("GET /stock/transactions", stockTransactionHandler.StockTransactionsPage)

	// Stock details page
	mux.HandleFunc("GET /stock/{id}", stockTransactionHandler.StockDetailsPage)

	// Post transaction pages
	mux.HandleFunc("GET /stock/post-transaction", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/stock/post-transaction/stock-movement", http.StatusSeeOther)
	})

	// Stock Movement
	mux.HandleFunc("GET /stock/post-transaction/stock-movement", stockTransactionHandler.PostStockMovementPage)
	mux.HandleFunc("POST /stock/post-transaction/stock-movement", stockTransactionHandler.PostStockMovement)

	// Production
	mux.HandleFunc("GET /stock/post-transaction/production", stockTransactionHandler.PostProductionPage)
	mux.HandleFunc("POST /stock/post-transaction/production", stockTransactionHandler.PostProduction)

	// Production Reversal
	mux.HandleFunc("GET /stock/post-transaction/production-reversal", stockTransactionHandler.PostProductionReversalPage)
	mux.HandleFunc("POST /stock/post-transaction/production-reversal", stockTransactionHandler.PostProductionReversal)

	// Consumption
	mux.HandleFunc("GET /stock/post-transaction/consumption", stockTransactionHandler.PostConsumptionPage)
	mux.HandleFunc("POST /stock/post-transaction/consumption", stockTransactionHandler.PostConsumption)

	// Consumption Reversal
	mux.HandleFunc("GET /stock/post-transaction/consumption-reversal", stockTransactionHandler.PostConsumptionReversalPage)
	mux.HandleFunc("POST /stock/post-transaction/consumption-reversal", stockTransactionHandler.PostConsumptionReversal)

}
