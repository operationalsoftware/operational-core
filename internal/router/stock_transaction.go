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

	mux.HandleFunc("GET /stock/test", stockTrxHandler.GetStockLevels)
	mux.HandleFunc("GET /stock/levels", stockTrxHandler.GetStockLevels)
	mux.HandleFunc("POST /stock/transaction", stockTrxHandler.CreateStockTransaction)
	mux.HandleFunc("GET /stock/transactions", stockTrxHandler.GetStockTransactions)

	// QRcode login page
	// mux.HandleFunc("GET /auth/password/qrcode", stockTrxHandler.QRcodeLogInPage)
	// mux.HandleFunc("POST /auth/password/qrcode", stockTrxHandler.QRcodeLogIn)

	// mux.HandleFunc("/auth/logout", stockTrxHandler.Logout)
}
