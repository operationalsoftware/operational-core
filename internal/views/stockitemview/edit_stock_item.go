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

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
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
			components.InlineStyle("/internal/views/userview/edit_user_page.css"),
			components.InlineStyle("/internal/views/stockitemview/add_stock_item.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/stockitemview/add_stock_item.js"),
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

	productTypeLabel := "Product Type"
	productTypeKey := "ProductType"
	productTypeValue := p.stockItem.ProductType
	productTypeError := ""
	if p.isSubmission || productTypeValue != "" {
		productTypeError = p.validationErrors.GetError(productTypeKey, productTypeLabel)
	}
	productTypeHelperType := components.InputHelperTypeNone
	if productTypeError != "" {
		productTypeHelperType = components.InputHelperTypeError
	}

	yarnTypeLabel := "Yarn Type"
	yarnTypeKey := "YarnType"
	yarnTypeValue := p.stockItem.YarnType
	yarnTypeError := ""
	if p.isSubmission || yarnTypeValue != "" {
		yarnTypeError = p.validationErrors.GetError(yarnTypeKey, yarnTypeLabel)
	}
	yarnTypeHelperType := components.InputHelperTypeNone
	if yarnTypeError != "" {
		yarnTypeHelperType = components.InputHelperTypeError
	}

	styleNumberLabel := "Style Number"
	styleNumberKey := "StyleNumber"
	styleNumberValue := p.stockItem.StyleNumber
	styleNumberError := ""
	if p.isSubmission || styleNumberValue != "" {
		styleNumberError = p.validationErrors.GetError(styleNumberKey, styleNumberLabel)
	}
	styleNumberHelperType := components.InputHelperTypeNone
	if styleNumberError != "" {
		styleNumberHelperType = components.InputHelperTypeError
	}

	colourLabel := "Colour"
	colourKey := "Colour"
	colourValue := p.stockItem.Colour
	colourError := ""
	if p.isSubmission || colourValue != "" {
		colourError = p.validationErrors.GetError(colourKey, colourLabel)
	}
	colourHelperType := components.InputHelperTypeNone
	if colourError != "" {
		colourHelperType = components.InputHelperTypeError
	}

	toeClosingLabel := "Toe Closing"
	toeClosingKey := "ToeClosing"
	toeClosingValue := p.stockItem.ToeClosing
	toeClosingError := ""
	if p.isSubmission || toeClosingValue != "" {
		toeClosingError = p.validationErrors.GetError(toeClosingKey, toeClosingLabel)
	}
	toeClosingHelperType := components.InputHelperTypeNone
	if toeClosingError != "" {
		toeClosingHelperType = components.InputHelperTypeError
	}

	sizeLabel := "Size"
	sizeKey := "Size"
	sizeValue := p.stockItem.Size
	sizeError := ""
	if p.isSubmission || sizeValue != "" {
		sizeError = p.validationErrors.GetError(sizeKey, sizeLabel)
	}
	sizeHelperType := components.InputHelperTypeNone
	if sizeError != "" {
		sizeHelperType = components.InputHelperTypeError
	}

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

	return h.Div(
		components.Form(
			h.ID("sku-generate-form"),
			h.Method("POST"),

			h.Div(
				h.Class("sku-grid"),

				h.Label(
					h.Class("sku-option"),
					h.P(
						g.Text("Product Type"),
					),
					components.SearchSelect(&components.SearchSelectProps{
						Name:        "ProductType",
						Mode:        "single",
						Placeholder: "Select product type",
						Options:     ProductTypes,
						Selected:    productTypeValue,
					}),
					g.If(
						productTypeError != "",
						components.InputHelper(&components.InputHelperProps{
							Label: productTypeError,
							Type:  productTypeHelperType,
						},
						),
					),
				),
				h.Label(
					h.Class("sku-option"),
					h.P(
						g.Text("Yarn Type"),
					),
					components.SearchSelect(&components.SearchSelectProps{
						Name:        "YarnType",
						Mode:        "single",
						Placeholder: "Select yarn type",
						Options:     YarnTypes,
						Selected:    yarnTypeValue,
					}),
					g.If(
						yarnTypeError != "",
						components.InputHelper(&components.InputHelperProps{
							Label: yarnTypeError,
							Type:  yarnTypeHelperType,
						},
						),
					),
				),
				h.Label(
					h.Class("sku-option"),
					h.P(
						g.Text("Style Number"),
					),
					components.SearchSelect(&components.SearchSelectProps{
						Name:        "StyleNumber",
						Mode:        "single",
						Placeholder: "Select style",
						Options:     StyleNumbers,
						Selected:    styleNumberValue,
					}),
					g.If(
						styleNumberError != "",
						components.InputHelper(&components.InputHelperProps{
							Label: styleNumberError,
							Type:  styleNumberHelperType,
						},
						),
					),
				),
				h.Label(
					h.Class("sku-option"),
					h.P(
						g.Text("Colour"),
					),
					components.SearchSelect(&components.SearchSelectProps{
						Name:        "Colour",
						Mode:        "single",
						Placeholder: "Select colour",
						Options:     Colours,
						Selected:    colourValue,
					}),
					g.If(
						colourError != "",
						components.InputHelper(&components.InputHelperProps{
							Label: colourError,
							Type:  colourHelperType,
						},
						),
					),
				),
				h.Label(
					h.Class("sku-option"),
					h.P(
						g.Text("Toe Closing"),
					),
					components.SearchSelect(&components.SearchSelectProps{
						Name:        "ToeClosing",
						Mode:        "single",
						Placeholder: "Select toe closing",
						Options:     ToeClosings,
						Selected:    toeClosingValue,
					}),
					g.If(
						toeClosingError != "",
						components.InputHelper(&components.InputHelperProps{
							Label: toeClosingError,
							Type:  toeClosingHelperType,
						},
						),
					),
				),
				h.Label(
					h.Class("sku-option"),
					h.P(
						g.Text("Size"),
					),
					components.SearchSelect(&components.SearchSelectProps{
						Name:        "Size",
						Mode:        "single",
						Placeholder: "Select size",
						Options:     Sizes,
						Selected:    sizeValue,
					}),
					g.If(
						sizeError != "",
						components.InputHelper(&components.InputHelperProps{
							Label: sizeError,
							Type:  sizeHelperType,
						},
						),
					),
				),
			),

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
						h.ReadOnly(),
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
		),
	)

}
