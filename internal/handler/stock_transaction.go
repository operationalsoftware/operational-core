package handler

import (
	"app/internal/model"
	"app/internal/service"
	"app/pkg/reqcontext"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

type StockTrxHandler struct {
	stockTrxService service.StockTrxService
}

func NewStockTrxHandler(stockTrxService service.StockTrxService) *StockTrxHandler {
	return &StockTrxHandler{stockTrxService: stockTrxService}
}

func (h *StockTrxHandler) CreateStockTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	lotNumber := "LOT123"

	stockTrxs := model.PostStockTransactionsInput{
		{
			Timestamp:     time.Now(),
			StockCode:     "STK001",
			Qty:           decimal.NewFromFloat32(12.0),
			FromAccount:   "ACC001",
			FromLocation:  "LOC001",
			FromBin:       "BIN001",
			FromLotNumber: &lotNumber,
			ToAccount:     "ACC002",
			ToLocation:    "LOC002",
			ToBin:         "BIN002",
			ToLotNumber:   &lotNumber,
		},
	}

	err := h.stockTrxService.PostStockTransaction(r.Context(), &stockTrxs, ctx.User.UserID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to create stock transaction", http.StatusInternalServerError)
		return
	}

	// err := json.Unmarshal(r.Body, &stockTrxInput)
	// if err != nil {
	// 	http.Error(w, "Error decoding url values", http.StatusBadRequest)
	// 	return
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Stock Transaction created successfully"})

}

func ptr[T any](v T) *T {
	return &v
}

func (h *StockTrxHandler) GetStockTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	// filterInput := model.GetTransactionsInput{
	// 	// Account:      ptr("INV-001"),
	// 	// StockCode:    ptr("STK-123"),
	// 	// Location:     ptr("LOC-A"),
	// 	// Bin:          ptr("BIN-5"),
	// 	// LotNumber:    ptr("LOT-77"),
	// 	// LTETimestamp: &now,
	// 	// Page:         1,
	// 	// PageSize:     20,
	// }

	input := model.GetTransactionsInput{
		StockCode: ptr("STK001"),
		Location:  ptr("LOC002"),
		Page:      1,
		PageSize:  10,
	}

	transactions, err := h.stockTrxService.GetStockTransaction(r.Context(), &input, ctx.User.UserID)
	if err != nil {
		http.Error(w, "Failed to get stock transactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": transactions})
}

func (h *StockTrxHandler) GetStockLevels(w http.ResponseWriter, r *http.Request) {
	// ctx := reqcontext.GetContext(r)

	input := &model.GetStockLevelsInput{
		// Account: sql.NullString{String: "ACC001", Valid: true},
	}

	levels, err := h.stockTrxService.GetStockLevels(r.Context(), input)
	if err != nil {
		http.Error(w, "Failed to GET stock levels", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": levels})
}
