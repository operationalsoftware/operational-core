package stockitemview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"net/url"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

type AddStockItemPageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
}

func AddStockItemPage(p *AddStockItemPageProps) g.Node {

	content := g.Group([]g.Node{

		addStockItemForm(&addStockItemFormProps{
			values:           p.Values,
			validationErrors: p.ValidationErrors,
			isSubmission:     p.IsSubmission,
		}),
	})

	return layout.Page(layout.PageProps{
		Ctx:   p.Ctx,
		Title: "Add New Stock Item",
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "package-variant-closed",
				Title:          "Stock Items",
				URLPart:        "stock-items",
			},
			{
				IconIdentifier: "plus",
				Title:          "Add",
			},
		},
		Content: content,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/stockitemview/add_page.css"),
		},
	})
}

type addStockItemFormProps struct {
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
}

func addStockItemForm(p *addStockItemFormProps) g.Node {

	stockCodeLabel := "Stock Code (SKU)"
	stockCodeKey := "StockCode"
	stockCodeValue := p.values.Get(stockCodeKey)
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
	descriptionValue := p.values.Get(descriptionKey)
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
			g.Text("Add Stock Item"),
		),
	)
}
