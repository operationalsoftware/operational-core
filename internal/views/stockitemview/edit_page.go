package stockitemview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/internal/model"
	"app/pkg/nilsafe"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"fmt"
	"net/url"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type EditStockItemPageProps struct {
	Ctx              reqcontext.ReqContext
	StockItem        model.StockItem
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
}

func EditStockItemPage(p *EditStockItemPageProps) g.Node {

	content := g.Group([]g.Node{
		editUserForm(&editUserFormProps{
			stockItem:        p.StockItem,
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
		}),
	})

	return layout.Page(layout.PageProps{
		Title: fmt.Sprintf("Edit Stock Item: %s", p.StockItem.StockCode),
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "package-variant-closed",
				Title:          "Stock Items",
				URLPart:        "stock-items",
			},
			{
				Title:   p.StockItem.StockCode,
				URLPart: p.StockItem.StockCode,
			},
			{
				IconIdentifier: "pencil",
				Title:          "Edit",
			},
		},
		Content: content,
		Ctx:     p.Ctx,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/stockitemview/edit_page.css"),
		},
	})
}

type editUserFormProps struct {
	stockItem        model.StockItem
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
}

// same as addUserForm, but no password fields
func editUserForm(p *editUserFormProps) g.Node {

	stockCodeLabel := "Stock Code (SKU)"
	stockCodeKey := "StockCode"
	var stockCodeValue string
	if p.values.Get(stockCodeKey) != "" {
		stockCodeValue = p.values.Get(stockCodeKey)
	} else {
		stockCodeValue = nilsafe.Str(&p.stockItem.StockCode)
	}
	stockCodeError := ""
	if p.isSubmission || stockCodeValue != "" {
		stockCodeError = p.validationErrors.GetError(stockCodeKey, stockCodeLabel)
	}
	stockCodeHelperType := components.InputHelperTypeNone
	if stockCodeError != "" {
		stockCodeHelperType = components.InputHelperTypeError
	}

	descriptionLabel := "Description"
	descriptionKey := "Description"
	var descriptionValue string
	if p.values.Get(descriptionKey) != "" {
		descriptionValue = p.values.Get(descriptionKey)
	} else {
		descriptionValue = nilsafe.Str(&p.stockItem.Description)
	}
	descriptionError := ""
	if p.isSubmission || descriptionValue != "" {
		descriptionError = p.validationErrors.GetError(descriptionKey, descriptionLabel)
	}
	descriptionHelperType := components.InputHelperTypeNone
	if descriptionError != "" {
		descriptionHelperType = components.InputHelperTypeError
	}

	return components.Form(
		h.ID("sku-generate-form"),
		h.Method("POST"),

		h.Div(
			h.ID("sku-preview"),
			h.Class("sku-preview"),

			components.Input(&components.InputProps{
				Label:       stockCodeLabel,
				Name:        stockCodeKey,
				Placeholder: "Stock code",
				HelperText:  stockCodeError,
				HelperType:  stockCodeHelperType,
				InputProps: []g.Node{
					h.Value(stockCodeValue),
					h.AutoComplete("off"),
				},
			}),
		),

		components.Input(&components.InputProps{
			Label:       descriptionLabel,
			Name:        descriptionKey,
			Placeholder: "Enter description",
			HelperText:  descriptionError,
			HelperType:  descriptionHelperType,
			InputProps: []g.Node{
				h.Value(descriptionValue),
				h.AutoComplete("off"),
			},
		}),

		components.Button(
			&components.ButtonProps{},
			h.Type("submit"),
			g.Text("Save"),
		),
	)

}
