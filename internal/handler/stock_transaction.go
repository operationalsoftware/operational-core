package handler

import (
	"app/internal/model"
	"app/internal/service"
	"app/internal/views/stockview"
	"app/pkg/appurl"
	"app/pkg/nilsafe"
	"app/pkg/reqcontext"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/shopspring/decimal"
)

type StockTrxHandler struct {
	stockTrxService service.StockTrxService
}

func NewStockTrxHandler(stockTrxService service.StockTrxService) *StockTrxHandler {
	return &StockTrxHandler{stockTrxService: stockTrxService}
}

func normalise(uv *model.GetTransactionsInput) {
	uv.StockCode = strings.ToUpper(strings.TrimSpace(uv.StockCode))

	uv.Location = strings.ToUpper(strings.TrimSpace(uv.Location))

	uv.Bin = strings.ToUpper(strings.TrimSpace(uv.Bin))

	uv.LotNumber = strings.ToUpper(strings.TrimSpace(uv.LotNumber))

}

// func normalise(uv *model.GetStockLevelsInput) {
// 	if uv.StockCode != nil {
// 		stCode := strings.ToUpper(strings.TrimSpace(*uv.StockCode))
// 		uv.StockCode = &stCode
// 	}

// 	if uv.Location != nil {
// 		loc := strings.ToUpper(strings.TrimSpace(*uv.Location))
// 		uv.Location = &loc
// 	}

// 	if uv.Bin != nil {
// 		bin := strings.ToUpper(strings.TrimSpace(*uv.Bin))
// 		uv.Bin = &bin
// 	}

// 	if uv.LotNumber != nil {
// 		lotNo := strings.ToUpper(strings.TrimSpace(*uv.LotNumber))
// 		uv.LotNumber = &lotNo
// 	}

// }

// Pages
func (h *StockTrxHandler) StockLevelsPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	var uv model.GetTransactionsInput

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	normalise(&uv)

	if uv.Page == 0 {
		uv.Page = 1
	}
	if uv.PageSize == 0 {
		uv.PageSize = stockview.HomePageDefaultPageSize
	}

	stockLevels, err := h.stockTrxService.GetStockLevels(r.Context(), &model.GetStockLevelsInput{
		Account:      uv.Account,
		StockCode:    uv.StockCode,
		Location:     uv.Location,
		Bin:          uv.Bin,
		LotNumber:    uv.LotNumber,
		LTETimestamp: uv.LTETimestamp,
		Page:         uv.Page,
		PageSize:     uv.PageSize,
	})
	if err != nil {
		http.Error(w, "Error fetching stock levels", http.StatusInternalServerError)
		return
	}

	_ = stockview.StockLevelsPage(stockview.StockLevelsPageProps{
		Ctx:          ctx,
		StockLevels:  &stockLevels,
		Account:      uv.Account,
		StockCode:    uv.StockCode,
		Location:     uv.Location,
		Bin:          uv.Bin,
		LotNumber:    uv.LotNumber,
		LTETimestamp: uv.LTETimestamp,
		Page:         uv.Page,
		PageSize:     uv.PageSize,
		Total:        len(stockLevels),
	}).
		Render(w)

}

func (h *StockTrxHandler) StockDetailsPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	stockCode := r.PathValue("id")

	_ = stockview.StockDetailPage(stockview.StockDetailPageProps{
		Ctx:       ctx,
		StockCode: stockCode,
	}).
		Render(w)

}

func (h *StockTrxHandler) StockTransactionsPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	var uv model.GetTransactionsInput

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	normalise(&uv)

	if uv.Page == 0 {
		uv.Page = 1
	}
	if uv.PageSize == 0 {
		uv.PageSize = stockview.HomePageDefaultPageSize
	}

	stockTransactions, err := h.stockTrxService.GetStockTransactions(r.Context(), &model.GetTransactionsInput{
		Account:      uv.Account,
		StockCode:    uv.StockCode,
		Location:     uv.Location,
		Bin:          uv.Bin,
		LotNumber:    uv.LotNumber,
		LTETimestamp: uv.LTETimestamp,
		Page:         uv.Page,
		PageSize:     uv.PageSize,
	}, ctx.User.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error fetching stock transactions", http.StatusInternalServerError)
		return
	}

	_ = stockview.StockTransactionsPage(&stockview.StockTransactionsPageProps{
		Ctx:               ctx,
		StockTransactions: &stockTransactions,
		Account:           uv.Account,
		StockCode:         uv.StockCode,
		Location:          uv.Location,
		Bin:               uv.Bin,
		LotNumber:         uv.LotNumber,
		LTETimestamp:      uv.LTETimestamp,
		Page:              uv.Page,
		PageSize:          uv.PageSize,
		Total:             len(stockTransactions),
	}).Render(w)

}

func (h *StockTrxHandler) PostStockMovementPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	transactionType := "Stock Movement"

	// perms := ctx.User.Permissions
	hasPermission := true
	// hasPermission := perms.Production.Admin || perms.SupplyChain.Admin

	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	type postMovementUrlValues struct {
		StockCode    string
		LotNumber    string
		Qty          decimal.Decimal
		FromLocation string
		FromBin      string
		ToLocation   string
		ToBin        string
		ReturnTo     *string
	}

	var uv postMovementUrlValues

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	_ = stockview.PostStockMovementPage(
		&stockview.PostStockMovementPageProps{
			Ctx:             ctx,
			StockCode:       uv.StockCode,
			LotNumber:       uv.LotNumber,
			Qty:             uv.Qty,
			FromLocation:    uv.FromLocation,
			FromBin:         uv.FromBin,
			ToLocation:      uv.ToLocation,
			ToBin:           uv.ToBin,
			ReturnTo:        uv.ReturnTo,
			TransactionType: transactionType,
		},
	).Render(w)

}

func (h *StockTrxHandler) PostProductionPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	transactionType := "Production"
	perms := ctx.User.Permissions
	hasPermission := perms.Production.Admin || perms.SupplyChain.Admin || true
	// hasPermission := perms.Production.Admin || perms.SupplyChain.Admin

	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_ = stockview.PostProductionPage(
		&stockview.PostGenericPageProps{
			Ctx:             ctx,
			TransactionType: transactionType,
		},
	).Render(w)

}

func (h *StockTrxHandler) PostProductionReversalPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	transactionType := "Production Reversal"

	perms := ctx.User.Permissions
	hasPermission := perms.Production.Admin || perms.SupplyChain.Admin || true
	// hasPermission := perms.Production.Admin || perms.SupplyChain.Admin

	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_ = stockview.PostProductionReversalPage(
		&stockview.PostGenericPageProps{
			Ctx:             ctx,
			TransactionType: transactionType,
		},
	).Render(w)

}

func (h *StockTrxHandler) PostConsumptionPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	transactionType := "Consumption"

	perms := ctx.User.Permissions
	hasPermission := perms.Production.Admin || perms.SupplyChain.Admin || true
	// hasPermission := perms.Production.Admin || perms.SupplyChain.Admin

	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_ = stockview.PostConsumptionPage(
		&stockview.PostGenericPageProps{
			Ctx:             ctx,
			TransactionType: transactionType,
		},
	).Render(w)

}

func (h *StockTrxHandler) PostConsumptionReversalPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	transactionType := "Consumption Reversal"

	perms := ctx.User.Permissions
	hasPermission := perms.Production.Admin || perms.SupplyChain.Admin || true
	// hasPermission := perms.Production.Admin || perms.SupplyChain.Admin

	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_ = stockview.PostConsumptionReversalPage(
		&stockview.PostGenericPageProps{
			Ctx:             ctx,
			TransactionType: transactionType,
		},
	).Render(w)

}

// Request handlers
func (h *StockTrxHandler) PostStockTransactions(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	lotNumber := "LOT123"

	stockTrxs := model.PostStockTransactionsInput{
		{
			Timestamp:     nil,
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
		// StockCode: ptr("STK001"),
		// Location:  ptr("LOC002"),
		Page:     1,
		PageSize: 10,
	}

	transactions, err := h.stockTrxService.GetStockTransactions(r.Context(), &input, ctx.User.UserID)
	if err != nil {
		http.Error(w, "Failed to get stock transactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": transactions})
}

func (h *StockTrxHandler) GetStockLevels(w http.ResponseWriter, r *http.Request) {

	input := &model.GetStockLevelsInput{}

	levels, err := h.stockTrxService.GetStockLevels(r.Context(), input)
	if err != nil {
		http.Error(w, "Failed to GET stock levels", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"data": levels})
}

// Post Stock Transaction handlers

func (h *StockTrxHandler) PostStockMovement(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	// perms := ctx.User.Permissions
	hasPermission := true
	transactionType := "Stock Movement"

	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var fd postStockMovementFormData

	err = appurl.Unmarshal(r.Form, &fd)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd = mouldPostStockMovementFormData(fd)

	renderWithError := func(errorText string) {
		_ = stockview.PostStockMovementPage(
			&stockview.PostStockMovementPageProps{
				Ctx:             ctx,
				StockCode:       fd.StockCode,
				LotNumber:       fd.LotNumber,
				Qty:             fd.Qty,
				FromLocation:    fd.FromLocation,
				FromBin:         fd.FromBin,
				ToLocation:      fd.ToLocation,
				ToBin:           fd.ToBin,
				ReturnTo:        fd.ReturnTo,
				ErrorText:       errorText,
				TransactionType: transactionType,
			},
		).Render(w)
	}

	if fd.Qty.LessThanOrEqual(decimal.Zero) {
		renderWithError("Qty must be greater than 0")
		return
	}

	if fd.StockCode == "" {
		renderWithError("Stock code cannot be empty")
		return
	}
	if fd.FromLocation == "" {
		renderWithError("From location cannot be empty")
		return
	}
	if fd.ToLocation == "" {
		renderWithError("To location cannot be empty")
		return
	}

	err = h.stockTrxService.PostStockTransaction(
		r.Context(),
		&model.PostStockTransactionsInput{{
			StockTransactionType: transactionType,
			StockCode:            fd.StockCode,
			Qty:                  fd.Qty,
			FromAccount:          "STOCK",
			FromLocation:         fd.FromLocation,
			FromBin:              fd.FromBin,
			FromLotNumber:        &fd.LotNumber,
			ToAccount:            "STOCK",
			ToLocation:           fd.ToLocation,
			ToBin:                fd.ToBin,
			ToLotNumber:          &fd.LotNumber,
			StockTransactionNote: fd.TransactionNote,
		}},
		ctx.User.UserID,
	)
	if err != nil {
		renderWithError(fmt.Sprintf("Error posting movement: %v", err))
		return
	}

	// success if we got here
	if fd.ReturnTo != nil {
		http.Redirect(w, r, nilsafe.Str(fd.ReturnTo), http.StatusFound)
	}

	_ = stockview.PostStockMovementPage(
		&stockview.PostStockMovementPageProps{
			Ctx:             ctx,
			SuccessText:     "Successfully posted stock movement",
			TransactionType: transactionType,
		},
	).Render(w)
}

func (h *StockTrxHandler) PostProduction(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	transactionType := "Production"
	// perms := ctx.User.Permissions
	hasPermission := true

	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var fd postFormData

	err = appurl.Unmarshal(r.Form, &fd)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd = mouldPostProductionFormData(fd)

	renderWithError := func(errorText string) {
		_ = stockview.PostProductionPage(
			&stockview.PostGenericPageProps{
				Ctx:             ctx,
				StockCode:       fd.StockCode,
				Location:        fd.Location,
				Bin:             fd.Bin,
				LotNumber:       fd.LotNumber,
				Qty:             fd.Qty,
				ErrorText:       errorText,
				TransactionType: transactionType,
			},
		).Render(w)
	}

	if fd.Qty.LessThanOrEqual(decimal.Zero) {
		renderWithError("Qty must be greater than 0")
		return
	}

	if fd.StockCode == "" {
		renderWithError("Stock code cannot be empty")
		return
	}
	if fd.Location == "" {
		renderWithError("Location cannot be empty")
		return
	}

	err = h.stockTrxService.PostStockTransaction(
		r.Context(),
		&model.PostStockTransactionsInput{{
			StockTransactionType: transactionType,
			StockCode:            fd.StockCode,
			Qty:                  fd.Qty,
			FromAccount:          "PRODUCTION",
			FromLocation:         fd.Location,
			FromBin:              fd.Bin,
			FromLotNumber:        &fd.LotNumber,
			ToAccount:            "STOCK",
			ToLocation:           fd.Location,
			ToBin:                fd.Bin,
			ToLotNumber:          &fd.LotNumber,
			StockTransactionNote: fd.TransactionNote,
		}},
		ctx.User.UserID,
	)

	if err != nil {
		renderWithError(err.Error())
		return
	}

	_ = stockview.PostProductionPage(
		&stockview.PostGenericPageProps{
			Ctx:             ctx,
			SuccessText:     "Production operation posted successfully",
			TransactionType: transactionType,
		},
	).Render(w)

}

func (h *StockTrxHandler) PostProductionReversal(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	transactionType := "Production Reversal"

	// perms := ctx.User.Permissions
	hasPermission := true

	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var fd postFormData

	err = appurl.Unmarshal(r.Form, &fd)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd = mouldPostProductionReversalFormData(fd)

	renderWithError := func(errorText string) {
		_ = stockview.PostProductionReversalPage(
			&stockview.PostGenericPageProps{
				Ctx:             ctx,
				StockCode:       fd.StockCode,
				Location:        fd.Location,
				Bin:             fd.Bin,
				LotNumber:       fd.LotNumber,
				Qty:             fd.Qty,
				ErrorText:       errorText,
				TransactionType: transactionType,
			},
		).Render(w)
	}

	if fd.Qty.LessThanOrEqual(decimal.Zero) {
		renderWithError("Qty must be greater than 0")
		return
	}

	if fd.StockCode == "" {
		renderWithError("Stock code cannot be empty")
		return
	}
	if fd.Location == "" {
		renderWithError("Location cannot be empty")
		return
	}

	err = h.stockTrxService.PostStockTransaction(
		r.Context(),
		&model.PostStockTransactionsInput{{
			StockTransactionType: transactionType,
			StockCode:            fd.StockCode,
			Qty:                  fd.Qty,
			FromAccount:          "STOCK",
			FromLocation:         fd.Location,
			FromBin:              fd.Bin,
			FromLotNumber:        &fd.LotNumber,
			ToAccount:            "PRODUCTION",
			ToLocation:           fd.Location,
			ToBin:                fd.Bin,
			ToLotNumber:          &fd.LotNumber,
			StockTransactionNote: fd.TransactionNote,
		}},
		ctx.User.UserID,
	)

	if err != nil {
		renderWithError(err.Error())
		return
	}

	_ = stockview.PostProductionReversalPage(
		&stockview.PostGenericPageProps{
			Ctx:             ctx,
			SuccessText:     "Production reversal operation posted successfully",
			TransactionType: transactionType,
		},
	).Render(w)
}

func (h *StockTrxHandler) PostConsumption(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	transactionType := "Consumption"

	// perms := ctx.User.Permissions
	hasPermission := true

	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var fd postFormData

	err = appurl.Unmarshal(r.Form, &fd)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd = mouldPostConsumptionFormData(fd)

	renderWithError := func(errorText string) {
		_ = stockview.PostConsumptionPage(
			&stockview.PostGenericPageProps{
				Ctx:             ctx,
				StockCode:       fd.StockCode,
				Location:        fd.Location,
				Bin:             fd.Bin,
				LotNumber:       fd.LotNumber,
				Qty:             fd.Qty,
				ErrorText:       errorText,
				TransactionType: transactionType,
			},
		).Render(w)
	}

	if fd.Qty.LessThanOrEqual(decimal.Zero) {
		renderWithError("Qty must be greater than 0")
		return
	}

	if fd.StockCode == "" {
		renderWithError("Stock code cannot be empty")
		return
	}
	if fd.Location == "" {
		renderWithError("Location cannot be empty")
		return
	}

	err = h.stockTrxService.PostStockTransaction(
		r.Context(),
		&model.PostStockTransactionsInput{{
			StockTransactionType: transactionType,
			StockCode:            fd.StockCode,
			Qty:                  fd.Qty,
			FromAccount:          "STOCK",
			FromLocation:         fd.Location,
			FromBin:              fd.Bin,
			FromLotNumber:        &fd.LotNumber,
			ToAccount:            "CONSUMED",
			ToLocation:           fd.Location,
			ToBin:                fd.Bin,
			ToLotNumber:          &fd.LotNumber,
			StockTransactionNote: fd.TransactionNote,
		}},
		ctx.User.UserID,
	)

	if err != nil {
		renderWithError(err.Error())
		return
	}

	_ = stockview.PostConsumptionPage(
		&stockview.PostGenericPageProps{
			Ctx:             ctx,
			SuccessText:     "Consumption operation posted successfully",
			TransactionType: transactionType,
		},
	).Render(w)
}

func (h *StockTrxHandler) PostConsumptionReversal(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	transactionType := "Consumption Reversal"

	// perms := ctx.User.Permissions
	hasPermission := true

	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var fd postFormData

	err = appurl.Unmarshal(r.Form, &fd)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd = mouldPostConsumptionReversalFormData(fd)

	renderWithError := func(errorText string) {
		_ = stockview.PostConsumptionReversalPage(
			&stockview.PostGenericPageProps{
				Ctx:             ctx,
				StockCode:       fd.StockCode,
				Location:        fd.Location,
				Bin:             fd.Bin,
				LotNumber:       fd.LotNumber,
				Qty:             fd.Qty,
				ErrorText:       errorText,
				TransactionType: transactionType,
			},
		).Render(w)
	}

	if fd.Qty.LessThanOrEqual(decimal.Zero) {
		renderWithError("Qty must be greater than 0")
		return
	}

	if fd.StockCode == "" {
		renderWithError("Stock code cannot be empty")
		return
	}
	if fd.Location == "" {
		renderWithError("Location cannot be empty")
		return
	}

	err = h.stockTrxService.PostStockTransaction(
		r.Context(),
		&model.PostStockTransactionsInput{{
			StockTransactionType: transactionType,
			StockCode:            fd.StockCode,
			Qty:                  fd.Qty,
			FromAccount:          "CONSUMED",
			FromLocation:         fd.Location,
			FromBin:              fd.Bin,
			FromLotNumber:        &fd.LotNumber,
			ToAccount:            "STOCK",
			ToLocation:           fd.Location,
			ToBin:                fd.Bin,
			ToLotNumber:          &fd.LotNumber,
			StockTransactionNote: fd.TransactionNote,
		}},
		ctx.User.UserID,
	)

	if err != nil {
		renderWithError(err.Error())
		return
	}

	_ = stockview.PostConsumptionReversalPage(
		&stockview.PostGenericPageProps{
			Ctx:             ctx,
			SuccessText:     "Consumption reversal operation posted successfully",
			TransactionType: transactionType,
		},
	).Render(w)
}

// FORM DATA
type postFormData struct {
	StockCode       string
	Location        string
	Bin             string
	LotNumber       string
	Qty             decimal.Decimal
	TransactionNote string
}

// Stock movement

type postStockMovementFormData struct {
	StockCode       string
	Account         string
	LotNumber       string
	Qty             decimal.Decimal
	FromLocation    string
	FromBin         string
	ToLocation      string
	ToBin           string
	TransactionNote string
	ReturnTo        *string
}

func mouldPostStockMovementFormData(fd postStockMovementFormData) postStockMovementFormData {

	// trim and uppercase
	fd.StockCode = strings.ToUpper(strings.TrimSpace(fd.StockCode))

	// trim
	fd.FromLocation = strings.TrimSpace(fd.FromLocation)
	fd.FromBin = strings.TrimSpace(fd.FromBin)
	fd.ToLocation = strings.TrimSpace(fd.ToLocation)
	fd.ToBin = strings.TrimSpace(fd.ToBin)
	fd.LotNumber = strings.TrimSpace(fd.LotNumber)
	fd.TransactionNote = strings.TrimSpace(fd.TransactionNote)

	return fd
}

// Production
func mouldPostProductionFormData(fd postFormData) postFormData {

	// trim and uppercase
	fd.StockCode = strings.ToUpper(strings.TrimSpace(fd.StockCode))

	// trim
	fd.Location = strings.TrimSpace(fd.Location)
	fd.Bin = strings.TrimSpace(fd.Bin)

	return fd
}

// Production Reversal
func mouldPostProductionReversalFormData(fd postFormData) postFormData {

	// trim and uppercase
	fd.StockCode = strings.ToUpper(strings.TrimSpace(fd.StockCode))

	// trim
	fd.Location = strings.TrimSpace(fd.Location)
	fd.Bin = strings.TrimSpace(fd.Bin)
	fd.LotNumber = strings.TrimSpace(fd.LotNumber)

	return fd
}

// Consumption
func mouldPostConsumptionFormData(fd postFormData) postFormData {

	// trim and uppercase
	fd.StockCode = strings.ToUpper(strings.TrimSpace(fd.StockCode))

	// trim
	fd.Location = strings.TrimSpace(fd.Location)
	fd.Bin = strings.TrimSpace(fd.Bin)
	fd.LotNumber = strings.TrimSpace(fd.LotNumber)

	return fd
}

// Consumption Reversal
func mouldPostConsumptionReversalFormData(fd postFormData) postFormData {

	// trim and uppercase
	fd.StockCode = strings.ToUpper(strings.TrimSpace(fd.StockCode))

	// trim
	fd.Location = strings.TrimSpace(fd.Location)
	fd.Bin = strings.TrimSpace(fd.Bin)
	fd.LotNumber = strings.TrimSpace(fd.LotNumber)

	return fd
}
