package stockitemview

import (
	"app/internal/components"
	"app/internal/layout"
	"app/pkg/reqcontext"
	"app/pkg/validate"
	"net/url"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
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
			components.InlineStyle("/internal/views/userview/add_user_page.css"),
			components.InlineStyle("/internal/views/stockitemview/add_stock_item.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/stockitemview/add_stock_item.js"),
		},
	})
}

type addStockItemFormProps struct {
	values           url.Values
	validationErrors validate.ValidationErrors
	isSubmission     bool
}

func addStockItemForm(p *addStockItemFormProps) g.Node {

	productTypeLabel := "Product Type"
	productTypeKey := "ProductType"
	productTypeValue := p.values.Get(productTypeKey)
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
	yarnTypeValue := p.values.Get(yarnTypeKey)
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
	styleNumberValue := p.values.Get(styleNumberKey)
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
	colourValue := p.values.Get(colourKey)
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
	toeClosingValue := p.values.Get(toeClosingKey)
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
	sizeValue := p.values.Get(sizeKey)
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
				g.Text("Add Stock Item"),
			),
		),
	)
}

var ProductTypes = []components.SearchSelectOption{
	{Label: "MENS PACK", Value: "30"},
	{Label: "LADIES PACK", Value: "31"},
	{Label: "MENS INVISIBLE", Value: "60"},
	{Label: "LADIES INVISIBLE", Value: "61"},
	{Label: "MENS TRAINER", Value: "70"},
	{Label: "LADIES TRAINER", Value: "71"},
	{Label: "MENS REGULAR", Value: "80"},
	{Label: "LADIES REGULAR", Value: "81"},
	{Label: "MENS LONG", Value: "90"},
	{Label: "LADIES LONG", Value: "91"},
}

var YarnTypes = []components.SearchSelectOption{
	{Label: "CASHMERE (HG)", Value: "10"},
	{Label: "CASHMERE/SILK (FG)", Value: "14"},
	{Label: "CASHMERE COTTON (HG) (64% CASHMERE 36% COTTON)", Value: "15"},
	{Label: "COTTON/CASHMERE (FG) (64% COTTON 25% NYLON 11% CASHMERE)", Value: "20"},
	{Label: "GLENN LYON (HG)", Value: "25"},
	{Label: "WOOL (HG) (80% WOOL 20% COTTON)", Value: "30"},
	{Label: "100% BRITISH WOOL (HG)", Value: "31"},
	{Label: "DONEGAL (HG) (70% MERINO WOOL, 10% SILK, 10% CASH, 10% COTT)", Value: "33"},
	{Label: "WOOL (FG) (63% WOOL, 37% NYLON)", Value: "35"},
	{Label: "TOLEGANO WOOL (FG) (LABARUNUM)", Value: "38"},
	{Label: "COTTON (HG) (100% SOFT COTTON)", Value: "40"},
	{Label: "COTTON (FG) (75% COTTON, 25% NYLON)", Value: "45"},
	{Label: "CASHMERE COTTON (HG) (64% CASHMERE 36% COTTON)", Value: "46"},
	{Label: "CASHMERE COTTON (HG) (64% CASHMERE 36% COTTON)", Value: "55"},
	{Label: "CASHMERE COTTON (HG) (64% CASHMERE 36% COTTON)", Value: "60"},
}

var StyleNumbers = []components.SearchSelectOption{
	{Label: "One", Value: "123"},
	{Label: "Two", Value: "456"},
	{Label: "Three", Value: "789"},
}

var Colours = []components.SearchSelectOption{
	{Label: "Black", Value: "BK"},
	{Label: "Blue", Value: "BL"},
	{Label: "Gray", Value: "GY"},
	{Label: "Red", Value: "RD"},
	{Label: "Green", Value: "GN"},
	{Label: "Yellow", Value: "YL"},
}

var Sizes = []components.SearchSelectOption{
	{Label: "Small", Value: "S"},
	{Label: "Medium", Value: "M"},
	{Label: "Large", Value: "L"},
}

var ToeClosings = []components.SearchSelectOption{
	{Label: "Rosso", Value: "R"},
	{Label: "Handlinked", Value: "H"},
	{Label: "Lonati", Value: "L"},
}
