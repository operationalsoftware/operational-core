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
	sort.ParseQueryParam(uv.Sort, []string{
		"StockCode",
		"Description",
		"CreatedAt",
	})

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

	stockCode := r.PathValue("stockCode")

	stockItem, err := h.stockItemService.GetStockItem(r.Context(), stockCode)
	if err != nil {
		http.Error(w, "Error fetching Stock item", http.StatusInternalServerError)
		return
	}

	stockItemChanges, err := h.stockItemService.GetStockItemChanges(r.Context(), stockCode)
	if err != nil {
		http.Error(w, "Error fetching Stock item changes", http.StatusInternalServerError)
		return
	}

	if stockItem == nil {
		http.Error(w, "Stock item not found", http.StatusNotFound)
		return
	}

	png, err := qrcode.Encode(stockCode, qrcode.Medium, 256)
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
		ProductType: formData.ProductType,
		YarnType:    formData.YarnType,
		StyleNumber: formData.StyleNumber,
		Colour:      formData.Colour,
		ToeClosing:  formData.ToeClosing,
		Size:        formData.Size,
		StockCode:   formData.StockCode,
		Description: formData.Description,
	}, ctx.User.UserID)

	fmt.Println(err)
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

	stockCode := r.PathValue("stockCode")

	stockItem, err := h.stockItemService.GetStockItem(r.Context(), stockCode)
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

	stockCode := r.PathValue("stockCode")

	stockItem, err := h.stockItemService.GetStockItem(r.Context(), stockCode)
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

	validationErrors, err := h.stockItemService.UpdateStockItem(r.Context(), stockCode, &model.PostStockItem{
		ProductType: formData.ProductType,
		YarnType:    formData.YarnType,
		StyleNumber: formData.StyleNumber,
		Colour:      formData.Colour,
		ToeClosing:  formData.ToeClosing,
		Size:        formData.Size,
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

	http.Redirect(w, r, "/stock-items/"+formData.StockCode, http.StatusSeeOther)
}

func (h *StockItemHandler) GenerateSKU(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)
	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	type urlVals struct {
		ProductType string
		YarnType    string
		StyleNumber string
		Colour      string
		Size        string
		ToeClosing  string
	}
	var uv urlVals

	err := appurl.Unmarshal(r.URL.Query(), &uv)
	if err != nil {
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}

	sku := fmt.Sprintf("%s-%s-%s-%s-%s-%s",
		uv.ProductType,
		uv.YarnType,
		uv.StyleNumber,
		uv.Colour,
		uv.ToeClosing,
		uv.Size,
	)

	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		_ = stockitemview.StockCodePartial(stockitemview.SKUPartialProps{
			SKU: sku,
		}).Render(w)
		return
	}
}

func (h *StockItemHandler) SKUConfigPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	skuItems, err := h.stockItemService.GetSKUConfiguration(r.Context())
	if err != nil {
		_ = stockitemview.SKUItemsPage(&stockitemview.SKUItemsPageProps{
			Ctx: ctx,
		}).
			Render(w)
		return
	}

	_ = stockitemview.SKUItemsPage(&stockitemview.SKUItemsPageProps{
		Ctx:      ctx,
		SKUItems: skuItems,
	}).
		Render(w)

}

func (h *StockItemHandler) AddSKUItemPage(w http.ResponseWriter, r *http.Request) {
	ctx := reqcontext.GetContext(r)

	hasPermission := ctx.User.Permissions.UserAdmin.Access
	if !hasPermission {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_ = stockitemview.AddSKUItemPage(&stockitemview.AddSKUItemPageProps{
		Ctx: ctx,
	}).Render(w)

}

func (h *StockItemHandler) AddSKUItem(w http.ResponseWriter, r *http.Request) {
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

	var uv postSKUItemFormData

	err = appurl.Unmarshal(r.Form, &uv)
	if err != nil {
		http.Error(w, "Error decoding url values", http.StatusBadRequest)
		return
	}
	uv.normalise()

	validationErrors, err := h.stockItemService.CreateSKUConfigItem(r.Context(), &model.SKUConfigItem{
		SKUField: uv.SKUField,
		Label:    uv.Label,
		Code:     uv.Code,
	})
	if len(validationErrors) > 0 {
		_ = stockitemview.AddSKUItemPage(&stockitemview.AddSKUItemPageProps{
			Ctx:              ctx,
			Values:           r.Form,
			ValidationErrors: validationErrors,
			IsSubmission:     true,
		}).Render(w)
		return
	}

	if err != nil {
		_ = stockitemview.AddSKUItemPage(&stockitemview.AddSKUItemPageProps{
			Ctx:       ctx,
			Values:    r.Form,
			ErrorText: err.Error(),
		}).Render(w)
		return
	}

	_ = stockitemview.AddSKUItemPage(&stockitemview.AddSKUItemPageProps{
		Ctx:         ctx,
		SuccessText: "SKU config added successfully.",
	}).Render(w)
}

func (h *StockItemHandler) DeleteSKUItem(w http.ResponseWriter, r *http.Request) {

	skuField := r.PathValue("skuField")
	skuCode := r.PathValue("code")

	err := h.stockItemService.DeleteSKUConfigItem(r.Context(), skuField, skuCode)
	if err != nil {
		http.Error(w, "Failed to delete SKU config item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type postStockItemFormData struct {
	ProductType string
	YarnType    string
	StyleNumber string
	Colour      string
	Size        string
	ToeClosing  string
	StockCode   string
	Description string
}

func (fd *postStockItemFormData) normalise() {
	// trim and uppercase
	fd.StockCode = strings.ToUpper(strings.TrimSpace(fd.StockCode))

	// trim
	fd.Description = strings.TrimSpace(fd.Description)
}

type postSKUItemFormData struct {
	SKUField string
	Label    string
	Code     string
}

func (fd *postSKUItemFormData) normalise() {
	// trim and uppercase
	if fd.Label != "" {
		fd.Label = strings.ToUpper(strings.TrimSpace(fd.Label))
	}

	// trim
	if fd.Code != "" {
		fd.Code = strings.TrimSpace(fd.Code)
	}
}
