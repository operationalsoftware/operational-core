package handler

import (
	"app/internal/model"
	"app/internal/service"
	"app/internal/views/stockview"
	"app/pkg/appurl"
	"app/pkg/nilsafe"
	"app/pkg/reqcontext"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/shopspring/decimal"
)

type StockTransactionHandler struct {
	stockTransactionService service.StockTransactionService
}

func NewStockTransactionHandler(stockTransactionService service.StockTransactionService) *StockTransactionHandler {
	return &StockTransactionHandler{stockTransactionService: stockTransactionService}
}

type stockInputURLVals struct {
	Account      string
	StockCode    string
	Location     string
	Bin          string
	LotNumber    string
	LTETimestamp *time.Time
	Page         int
	PageSize     int
}

func (uv *stockInputURLVals) normalise() {

	if uv.Account == "" {
		uv.Account = "STOCK" // Default
	}

	uv.StockCode = strings.ToUpper(strings.TrimSpace(uv.StockCode))

	uv.Location = strings.ToUpper(strings.TrimSpace(uv.Location))

	uv.Bin = strings.ToUpper(strings.TrimSpace(uv.Bin))

	uv.LotNumber = strings.ToUpper(strings.TrimSpace(uv.LotNumber))

	if uv.Page == 0 {
		uv.Page = 1
	}
	if uv.PageSize == 0 {
		uv.PageSize = stockview.HomePageDefaultPageSize
	}
}

// Pages
func (h *StockTransactionHandler) StockLevelsPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	var uv stockInputURLVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	uv.normalise()

	stockLevels, err := h.stockTransactionService.GetStockLevels(r.Context(), &model.GetStockLevelsInput{
		Account:      model.StockAccount(uv.Account),
		StockCode:    uv.StockCode,
		Location:     uv.Location,
		Bin:          uv.Bin,
		LotNumber:    uv.LotNumber,
		LTETimestamp: uv.LTETimestamp,
		Page:         uv.Page,
		PageSize:     uv.PageSize,
	})
	if err != nil {
		log.Println(err)
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

func (h *StockTransactionHandler) StockDetailsPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	stockCode := r.PathValue("id")

	_ = stockview.StockDetailPage(stockview.StockDetailPageProps{
		Ctx:       ctx,
		StockCode: stockCode,
	}).
		Render(w)

}

func (h *StockTransactionHandler) StockTransactionsPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	var uv stockInputURLVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	uv.normalise()

	stockTransactions, err := h.stockTransactionService.GetStockTransactions(r.Context(), &model.GetTransactionsInput{
		Account:      model.StockAccount(uv.Account),
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

func (h *StockTransactionHandler) PostStockMovementPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	perms := ctx.User.Permissions
	hasPermission := perms.SupplyChain.Admin

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
			Ctx:          ctx,
			StockCode:    uv.StockCode,
			LotNumber:    uv.LotNumber,
			Qty:          uv.Qty,
			FromLocation: uv.FromLocation,
			FromBin:      uv.FromBin,
			ToLocation:   uv.ToLocation,
			ToBin:        uv.ToBin,
			ReturnTo:     uv.ReturnTo,
		},
	).Render(w)

}

func (h *StockTransactionHandler) PostProductionPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	perms := ctx.User.Permissions
	hasPermission := perms.SupplyChain.Admin

	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_ = stockview.PostProductionPage(
		&stockview.PostGenericPageProps{
			Ctx: ctx,
		},
	).Render(w)

}

func (h *StockTransactionHandler) PostProductionReversalPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	perms := ctx.User.Permissions
	hasPermission := perms.SupplyChain.Admin

	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_ = stockview.PostProductionReversalPage(
		&stockview.PostGenericPageProps{
			Ctx: ctx,
		},
	).Render(w)

}

func (h *StockTransactionHandler) PostConsumptionPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	perms := ctx.User.Permissions
	hasPermission := perms.SupplyChain.Admin

	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_ = stockview.PostConsumptionPage(
		&stockview.PostGenericPageProps{
			Ctx: ctx,
		},
	).Render(w)

}

func (h *StockTransactionHandler) PostConsumptionReversalPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	perms := ctx.User.Permissions
	hasPermission := perms.SupplyChain.Admin

	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_ = stockview.PostConsumptionReversalPage(
		&stockview.PostGenericPageProps{
			Ctx: ctx,
		},
	).Render(w)

}

func (h *StockTransactionHandler) PostStockAdjustmentPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	perms := ctx.User.Permissions
	hasPermission := perms.SupplyChain.Admin

	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_ = stockview.PostStockAdjustPage(
		&stockview.PostGenericPageProps{
			Ctx:               ctx,
			IsStockAdjustment: true,
		},
	).Render(w)

}

// Post Stock Transaction handlers

func (h *StockTransactionHandler) PostStockMovement(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	perms := ctx.User.Permissions
	hasPermission := perms.SupplyChain.Admin

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

	fd.normalise()

	renderWithError := func(errorText string) {
		_ = stockview.PostStockMovementPage(
			&stockview.PostStockMovementPageProps{
				Ctx:          ctx,
				StockCode:    fd.StockCode,
				LotNumber:    fd.LotNumber,
				Qty:          fd.Qty,
				FromLocation: fd.FromLocation,
				FromBin:      fd.FromBin,
				ToLocation:   fd.ToLocation,
				ToBin:        fd.ToBin,
				ReturnTo:     fd.ReturnTo,
				ErrorText:    errorText,
			},
		).Render(w)
	}

	errorText := fd.validate()
	if errorText != "" {
		renderWithError(errorText)
		return
	}

	err = h.stockTransactionService.PostManualStockMovement(
		r.Context(),
		&model.PostManualStockMovementInput{
			StockCode:       fd.StockCode,
			Qty:             fd.Qty,
			FromLocation:    fd.FromLocation,
			FromBin:         fd.FromBin,
			LotNumber:       fd.LotNumber,
			ToLocation:      fd.ToLocation,
			ToBin:           fd.ToBin,
			TransactionNote: fd.TransactionNote,
		},
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
			Ctx:         ctx,
			SuccessText: "Successfully posted stock movement",
		},
	).Render(w)
}

func (h *StockTransactionHandler) PostProduction(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	perms := ctx.User.Permissions
	hasPermission := perms.SupplyChain.Admin

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

	var fd postGenericTransactionFormData

	err = appurl.Unmarshal(r.Form, &fd)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd.normalise()

	renderWithError := func(errorText string) {
		_ = stockview.PostProductionPage(
			&stockview.PostGenericPageProps{
				Ctx:       ctx,
				StockCode: fd.StockCode,
				Location:  fd.Location,
				Bin:       fd.Bin,
				LotNumber: fd.LotNumber,
				Qty:       fd.Qty,
				ErrorText: errorText,
			},
		).Render(w)
	}

	errorText := fd.validate()
	if errorText != "" {
		renderWithError(errorText)
		return
	}

	err = h.stockTransactionService.PostManualProduction(
		r.Context(),
		&model.PostManualGenericStockTransactionInput{
			StockCode:       fd.StockCode,
			Qty:             fd.Qty,
			Location:        fd.Location,
			Bin:             fd.Bin,
			LotNumber:       fd.LotNumber,
			TransactionNote: fd.TransactionNote,
		},
		ctx.User.UserID,
	)

	if err != nil {
		renderWithError(err.Error())
		return
	}

	_ = stockview.PostProductionPage(
		&stockview.PostGenericPageProps{
			Ctx:         ctx,
			SuccessText: "Production operation posted successfully",
		},
	).Render(w)

}

func (h *StockTransactionHandler) PostProductionReversal(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	perms := ctx.User.Permissions
	hasPermission := perms.SupplyChain.Admin

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

	var fd postGenericTransactionFormData

	err = appurl.Unmarshal(r.Form, &fd)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd.normalise()

	renderWithError := func(errorText string) {
		_ = stockview.PostProductionReversalPage(
			&stockview.PostGenericPageProps{
				Ctx:       ctx,
				StockCode: fd.StockCode,
				Location:  fd.Location,
				Bin:       fd.Bin,
				LotNumber: fd.LotNumber,
				Qty:       fd.Qty,
				ErrorText: errorText,
			},
		).Render(w)
	}

	errorText := fd.validate()
	if errorText != "" {
		renderWithError(errorText)
		return
	}

	err = h.stockTransactionService.PostManualProductionReversal(
		r.Context(),
		&model.PostManualGenericStockTransactionInput{
			StockCode:       fd.StockCode,
			Qty:             fd.Qty,
			Location:        fd.Location,
			Bin:             fd.Bin,
			LotNumber:       fd.LotNumber,
			TransactionNote: fd.TransactionNote,
		},
		ctx.User.UserID,
	)

	if err != nil {
		renderWithError(err.Error())
		return
	}

	_ = stockview.PostProductionReversalPage(
		&stockview.PostGenericPageProps{
			Ctx:         ctx,
			SuccessText: "Production reversal operation posted successfully",
		},
	).Render(w)
}

func (h *StockTransactionHandler) PostConsumption(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	perms := ctx.User.Permissions
	hasPermission := perms.SupplyChain.Admin

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

	var fd postGenericTransactionFormData

	err = appurl.Unmarshal(r.Form, &fd)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd.normalise()

	renderWithError := func(errorText string) {
		_ = stockview.PostConsumptionPage(
			&stockview.PostGenericPageProps{
				Ctx:       ctx,
				StockCode: fd.StockCode,
				Location:  fd.Location,
				Bin:       fd.Bin,
				LotNumber: fd.LotNumber,
				Qty:       fd.Qty,
				ErrorText: errorText,
			},
		).Render(w)
	}

	errorText := fd.validate()
	if errorText != "" {
		renderWithError(errorText)
		return
	}

	err = h.stockTransactionService.PostManualConsumption(
		r.Context(),
		&model.PostManualGenericStockTransactionInput{
			StockCode:       fd.StockCode,
			Qty:             fd.Qty,
			Location:        fd.Location,
			Bin:             fd.Bin,
			LotNumber:       fd.LotNumber,
			TransactionNote: fd.TransactionNote,
		},
		ctx.User.UserID,
	)

	if err != nil {
		renderWithError(err.Error())
		return
	}

	_ = stockview.PostConsumptionPage(
		&stockview.PostGenericPageProps{
			Ctx:         ctx,
			SuccessText: "Consumption operation posted successfully",
		},
	).Render(w)
}

func (h *StockTransactionHandler) PostConsumptionReversal(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	perms := ctx.User.Permissions
	hasPermission := perms.SupplyChain.Admin

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

	var fd postGenericTransactionFormData

	err = appurl.Unmarshal(r.Form, &fd)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd.normalise()

	renderWithError := func(errorText string) {
		_ = stockview.PostConsumptionReversalPage(
			&stockview.PostGenericPageProps{
				Ctx:       ctx,
				StockCode: fd.StockCode,
				Location:  fd.Location,
				Bin:       fd.Bin,
				LotNumber: fd.LotNumber,
				Qty:       fd.Qty,
				ErrorText: errorText,
			},
		).Render(w)
	}

	errorText := fd.validate()
	if errorText != "" {
		renderWithError(errorText)
		return
	}

	err = h.stockTransactionService.PostManualConsumptionReversal(
		r.Context(),
		&model.PostManualGenericStockTransactionInput{
			StockCode:       fd.StockCode,
			Qty:             fd.Qty,
			Location:        fd.Location,
			Bin:             fd.Bin,
			LotNumber:       fd.LotNumber,
			TransactionNote: fd.TransactionNote,
		},
		ctx.User.UserID,
	)

	if err != nil {
		renderWithError(err.Error())
		return
	}

	_ = stockview.PostConsumptionReversalPage(
		&stockview.PostGenericPageProps{
			Ctx:         ctx,
			SuccessText: "Consumption reversal operation posted successfully",
		},
	).Render(w)
}

func (h *StockTransactionHandler) PostStockAdjustment(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	perms := ctx.User.Permissions
	hasPermission := perms.SupplyChain.Admin

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

	var fd postGenericTransactionFormData
	fd.IsStockAdjustment = true

	err = appurl.Unmarshal(r.Form, &fd)
	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	fd.normalise()

	renderWithError := func(errorText string) {
		_ = stockview.PostStockAdjustPage(
			&stockview.PostGenericPageProps{
				Ctx:       ctx,
				StockCode: fd.StockCode,
				Location:  fd.Location,
				Bin:       fd.Bin,
				LotNumber: fd.LotNumber,
				Qty:       fd.Qty,
				ErrorText: errorText,
			},
		).Render(w)
	}

	errorText := fd.validate()
	if errorText != "" {
		renderWithError(errorText)
		return
	}

	err = h.stockTransactionService.PostManualStockAdjustment(
		r.Context(),
		&model.PostManualGenericStockTransactionInput{
			StockCode:       fd.StockCode,
			Qty:             fd.Qty,
			Location:        fd.Location,
			Bin:             fd.Bin,
			LotNumber:       fd.LotNumber,
			TransactionNote: fd.TransactionNote,
		},
		ctx.User.UserID,
	)

	if err != nil {
		renderWithError(err.Error())
		return
	}

	_ = stockview.PostStockAdjustPage(
		&stockview.PostGenericPageProps{
			Ctx:         ctx,
			SuccessText: "Stock adjustment operation posted successfully",
		},
	).Render(w)
}

type postGenericTransactionFormData struct {
	StockCode         string
	Location          string
	Bin               string
	LotNumber         string
	Qty               decimal.Decimal
	TransactionNote   string
	IsStockAdjustment bool
}

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

func (fd *postStockMovementFormData) normalise() {

	// trim and uppercase
	fd.StockCode = strings.ToUpper(strings.TrimSpace(fd.StockCode))
	fd.FromLocation = strings.ToUpper(strings.TrimSpace(fd.FromLocation))
	fd.FromBin = strings.ToUpper(strings.TrimSpace(fd.FromBin))
	fd.ToLocation = strings.ToUpper(strings.TrimSpace(fd.ToLocation))
	fd.ToBin = strings.ToUpper(strings.TrimSpace(fd.ToBin))
	fd.LotNumber = strings.ToUpper(strings.TrimSpace(fd.LotNumber))

	// trim
	fd.TransactionNote = strings.TrimSpace(fd.TransactionNote)

}

func (fd *postStockMovementFormData) validate() string {

	if fd.Qty.LessThanOrEqual(decimal.Zero) {
		return "Qty must be greater than 0"
	}

	if fd.StockCode == "" {
		return "Stock code cannot be empty"
	}
	if fd.FromLocation == "" {
		return "From location cannot be empty"
	}
	if fd.ToLocation == "" {
		return "To location cannot be empty"
	}

	return ""
}

func (fd *postGenericTransactionFormData) normalise() {

	// trim and uppercase
	fd.StockCode = strings.ToUpper(strings.TrimSpace(fd.StockCode))
	fd.Location = strings.ToUpper(strings.TrimSpace(fd.Location))
	fd.Bin = strings.ToUpper(strings.TrimSpace(fd.Bin))
	fd.LotNumber = strings.ToUpper(strings.TrimSpace(fd.LotNumber))

	// trim
	fd.TransactionNote = strings.TrimSpace(fd.TransactionNote)

}

func (fd *postGenericTransactionFormData) validate() string {

	if fd.Qty.LessThanOrEqual(decimal.Zero) && !fd.IsStockAdjustment {
		return "Qty must be greater than 0"
	}

	if fd.StockCode == "" {
		return "Stock code cannot be empty"
	}
	if fd.Location == "" {
		return "Location cannot be empty"
	}

	return ""
}
