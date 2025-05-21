package handler

import (
	"app/internal/model"
	"app/internal/service"
	"app/pkg/reqcontext"
	"encoding/json"
	"fmt"
	"net/http"

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

	fromLot := "LOT123"
	toLot := "LOT456"

	// stockTrxs := model.PostStockTransactionsInput{}

	stockTrxs := model.PostStockTransactionsInput{
		{
			// Timestamp:     time.Now(),
			StockCode:     "STK001",
			Qty:           decimal.NewFromFloat32(12.0),
			FromAccount:   "ACC001",
			FromLocation:  "LOC001",
			FromBin:       "BIN001",
			FromLotNumber: &fromLot,
			ToAccount:     "ACC002",
			ToLocation:    "LOC002",
			ToBin:         "BIN002",
			ToLotNumber:   &toLot,
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

// func (h *StockTrxHandler) TestStock(w http.ResponseWriter, r *http.Request) {
// 	// ctx := reqcontext.GetContext(r)
// 	fmt.Println("TESTING 123.....")

// 	return
// }

func ptr[T any](v T) *T {
	return &v
}

func (h *StockTrxHandler) GetStockTransactions(w http.ResponseWriter, r *http.Request) {
	// ctx := reqcontext.GetContext(r)
	// now := time.Now()

	// filterInput := model.GetTransactionsInput{
	// 	Account:      ptr("INV-001"),
	// 	StockCode:    ptr("STK-123"),
	// 	Location:     ptr("LOC-A"),
	// 	Bin:          ptr("BIN-5"),
	// 	LotNumber:    ptr("LOT-77"),
	// 	LTETimestamp: &now,
	// 	Page:         1,
	// 	PageSize:     20,
	// }

	// input := model.GetTransactionsInput{
	// 	StockCode: ptr("STK-123"),
	// 	Page:      1,
	// 	PageSize:  10,
	// }

	// transactions, err := h.stockTrxService.GetStockTransaction(r.Context(), &stockTrxs, ctx.User.UserID)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Error(w, "Failed to create stock transaction", http.StatusInternalServerError)
	// 	return
	// }

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated)
	// json.NewEncoder(w).Encode(map[string]interface{}{"data": transactions})
}

func (h *StockTrxHandler) GetStockLevels(w http.ResponseWriter, r *http.Request) {
	// ctx := reqcontext.GetContext(r)

	// var params model.SearchInput

	// results, err := h.searchService.Search(r.Context(), params.Q, allowedSearchEntities, ctx.User.UserID)
	// if err != nil {
	// 	_ = searchview.SearchPage(searchview.SearchPageProps{
	// 		Ctx:             ctx,
	// 		SearchTerm:      params.Q,
	// 		SearchEntities:  allowedSearchEntities,
	// 		UserPermissions: ctx.User.Permissions,
	// 	}).
	// 		Render(w)
	// 	return
	// }

	// _ = searchview.SearchPage(searchview.SearchPageProps{
	// 	Ctx:             ctx,
	// 	SearchTerm:      params.Q,
	// 	SearchEntities:  allowedSearchEntities,
	// 	Results:         results,
	// 	UserPermissions: ctx.User.Permissions,
	// }).
	// 	Render(w)

	return
}
