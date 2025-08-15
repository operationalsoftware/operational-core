package handler

import (
	"app/internal/model"
	"app/internal/service"
	"app/internal/views/stockitemview"
	"app/pkg/appsort"
	"app/pkg/appurl"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/skip2/go-qrcode"
)

type StockItemHandler struct {
	stockItemService service.StockItemService
}

func NewStockItemHandler(stockItemService service.StockItemService) *StockItemHandler {
	return &StockItemHandler{stockItemService: stockItemService}
}

func (h *StockItemHandler) StockItemsPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	type urlVals struct {
		Sort     string
		Page     int
		PageSize int
	}

	var uv urlVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	sort := appsort.Sort{}
	sort.ParseQueryParam(model.StockItem{}, uv.Sort)

	if uv.Page == 0 {
		uv.Page = 1
	}
	if uv.PageSize == 0 {
		uv.PageSize = 50
	}

	stockItems, count, err := h.stockItemService.GetStockItems(r.Context(), &model.GetStockItemsQuery{
		Sort:     sort,
		Page:     uv.Page,
		PageSize: uv.PageSize,
	}, ctx.User.UserID)
	if err != nil {
		_ = stockitemview.StockItemsPage(&stockitemview.StockItemsPageProps{
			Ctx: ctx,
		}).
			Render(w)
		return
	}

	_ = stockitemview.StockItemsPage(&stockitemview.StockItemsPageProps{
		Ctx:             ctx,
		StockItems:      stockItems,
		StockItemsCount: count,
		Sort:            sort,
		Page:            uv.Page,
		PageSize:        uv.PageSize,
	}).
		Render(w)

}

func (h *StockItemHandler) StockItemDetailsPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	stockItemID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid stock item ID", http.StatusBadRequest)
		return
	}

	stockItem, err := h.stockItemService.GetStockItem(r.Context(), stockItemID)
	if err != nil {
		http.Error(w, "Error fetching Stock item", http.StatusInternalServerError)
		return
	}

	stockItemChanges, err := h.stockItemService.GetStockItemChanges(r.Context(), stockItemID)
	if err != nil {
		http.Error(w, "Error fetching Stock item changes", http.StatusInternalServerError)
		return
	}

	if stockItem == nil {
		http.Error(w, "Stock item not found", http.StatusNotFound)
		return
	}

	png, err := qrcode.Encode(stockItem.StockCode, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "QR generation failed", http.StatusInternalServerError)
		return
	}

	// Convert PNG bytes to base64 string
	base64Image := base64.StdEncoding.EncodeToString(png)

	qrCodeURI := fmt.Sprintf("data:image/png;base64,%s", base64Image)

	_ = stockitemview.StockItemDetailsPage(&stockitemview.StockItemDetailsPageProps{
		Ctx:              ctx,
		StockItem:        stockItem,
		QRCode:           qrCodeURI,
		StockItemChanges: stockItemChanges,
	}).
		Render(w)
}

func (h *StockItemHandler) AddStockItemPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	query := r.URL.Query()

	var formData postStockItemFormData
	err := appurl.Unmarshal(query, &formData)
	if err != nil {
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	_ = stockitemview.AddStockItemPage(&stockitemview.AddStockItemPageProps{
		Ctx: ctx,
	}).Render(w)
}

func (h *StockItemHandler) AddStockItem(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	var formData postStockItemFormData
	var validationErrors validate.ValidationErrors
	err = appurl.Unmarshal(r.Form, &formData)

	formData.normalise()

	if err != nil {
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	validationErrors, err = h.stockItemService.CreateStockItem(r.Context(), &model.PostStockItem{
		StockCode:   formData.StockCode,
		Description: formData.Description,
	}, ctx.User.UserID)
	if err != nil {
		http.Error(w, "Error adding Stock item", http.StatusInternalServerError)
		return
	}

	if len(validationErrors) > 0 {
		_ = stockitemview.AddStockItemPage(&stockitemview.AddStockItemPageProps{
			Ctx:              ctx,
			Values:           r.Form,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
		}).Render(w)
		return
	}

	http.Redirect(w, r, "/stock-items", http.StatusSeeOther)
}

func (h *StockItemHandler) EditStockItemPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	stockItemID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid stock item ID", http.StatusBadRequest)
		return
	}

	stockItem, err := h.stockItemService.GetStockItem(r.Context(), stockItemID)
	if err != nil {
		http.Error(w, "Error getting Stock item", http.StatusInternalServerError)
		return
	}

	if stockItem == nil {
		http.Error(w, "Stock item does not exist", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	_ = stockitemview.EditStockItemPage(&stockitemview.EditStockItemPageProps{
		Ctx:       ctx,
		StockItem: *stockItem,
		Values:    r.Form,
	}).Render(w)
}

func (h *StockItemHandler) EditStockItem(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	stockItemID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid stock item ID", http.StatusBadRequest)
		return
	}

	stockItem, err := h.stockItemService.GetStockItem(r.Context(), stockItemID)
	if err != nil {
		http.Error(w, "Error getting Stock item", http.StatusInternalServerError)
		return
	}

	if stockItem == nil {
		http.Error(w, "Stock item does not exist", http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	values := r.Form

	var formData postStockItemFormData
	err = appurl.Unmarshal(values, &formData)
	if err != nil {
		log.Println("form unmarshal error:", err)
		http.Error(w, "Error decoding form", http.StatusBadRequest)
		return
	}

	formData.normalise()

	validationErrors, err := h.stockItemService.UpdateStockItem(r.Context(), stockItemID, &model.PostStockItem{
		StockCode:   formData.StockCode,
		Description: formData.Description,
	}, ctx.User.UserID)

	if err != nil {
		log.Println("update error:", err)
		http.Error(w, "Error updating Stock item", http.StatusInternalServerError)
		return
	}

	if len(validationErrors) > 0 {
		_ = stockitemview.EditStockItemPage(&stockitemview.EditStockItemPageProps{
			Ctx:              ctx,
			StockItem:        *stockItem,
			Values:           values,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
		}).Render(w)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/stock-items/%d", stockItemID), http.StatusSeeOther)
}

type postStockItemFormData struct {
	StockCode   string
	Description string
}

func (fd *postStockItemFormData) normalise() {
	// trim and uppercase
	fd.StockCode = strings.ToUpper(strings.TrimSpace(fd.StockCode))

	// trim
	fd.Description = strings.TrimSpace(fd.Description)
}

func (h *StockItemHandler) GetStockCodes(w http.ResponseWriter, r *http.Request) {
	searchText := r.URL.Query().Get("SearchText")
	selectedRaw := r.URL.Query().Get("SelectedValues")

	var selectedIDs []int
	if selectedRaw != "" {
		for _, v := range strings.Split(selectedRaw, ",") {
			if id, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
				selectedIDs = append(selectedIDs, id)
			}
		}
	}

	if len(selectedIDs) == 0 {
		selectedIDs = []int{}
	}

	stockItems, err := h.stockItemService.GetStockCodes(r.Context(), searchText, selectedIDs)
	if err != nil {
		log.Println(err)
	}

	for _, opt := range stockItems {
		fmt.Fprintf(w, `<div class="select-option" data-value="%d">%s</div>`, opt.StockItemID, opt.StockCode)
	}

}
