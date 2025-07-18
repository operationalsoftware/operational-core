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

type AddSKUItemPageProps struct {
	Ctx              reqcontext.ReqContext
	Values           url.Values
	ValidationErrors validate.ValidationErrors
	IsSubmission     bool
	SuccessText      string
	ErrorText        string
}

func AddSKUItemPage(p *AddSKUItemPageProps) g.Node {

	skuFieldLabel := "SKU Field"
	skuFieldKey := "SKUField"
	skuFieldValue := p.Values.Get(skuFieldKey)
	skuFieldError := ""
	if p.IsSubmission || skuFieldValue != "" {
		skuFieldError = p.ValidationErrors.GetError(skuFieldKey, skuFieldLabel)
	}
	skuFieldHelperType := components.InputHelperTypeNone
	if skuFieldError != "" {
		skuFieldHelperType = components.InputHelperTypeError
	}

	skuLabelLabel := "Label"
	skuLabelKey := "Label"
	skuLabelValue := p.Values.Get(skuLabelKey)
	skuLabelError := ""
	if p.IsSubmission || skuLabelValue != "" {
		skuLabelError = p.ValidationErrors.GetError(skuLabelKey, skuLabelLabel)
	}
	skuLabelHelperType := components.InputHelperTypeNone
	if skuLabelError != "" {
		skuLabelHelperType = components.InputHelperTypeError
	}

	codeLabel := "Code"
	codeKey := "Code"
	codeValue := p.Values.Get(codeKey)
	codeError := ""
	if p.IsSubmission || codeValue != "" {
		codeError = p.ValidationErrors.GetError(codeKey, codeLabel)
	}
	codeHelperType := components.InputHelperTypeNone
	if codeError != "" {
		codeHelperType = components.InputHelperTypeError
	}

	content := g.Group([]g.Node{

		components.Form(
			h.ID("sku-config-form"),
			h.Method("POST"),

			h.Label(
				h.Class("sku-field"),
				h.P(
					g.Text("SKU Field"),
				),
				components.SearchSelect(&components.SearchSelectProps{
					Name:          skuFieldKey,
					Mode:          "single",
					Placeholder:   "Select SKU field",
					Options:       SKUFields,
					Selected:      skuFieldValue,
					ShowOnlyLabel: true,
				}),
				g.If(
					skuFieldError != "",
					components.InputHelper(&components.InputHelperProps{
						Label: skuFieldError,
						Type:  skuFieldHelperType,
					},
					),
				),
			),

			components.Input(&components.InputProps{
				Label:       skuLabelLabel,
				Name:        skuLabelKey,
				Placeholder: "Enter label",
				HelperText:  skuLabelError,
				HelperType:  skuLabelHelperType,
				InputProps: []g.Node{
					h.Value(skuLabelValue),
					h.AutoComplete("off"),
				},
			}),

			h.Div(
				components.Input(&components.InputProps{
					Label:       codeLabel,
					Name:        codeKey,
					Placeholder: "Enter code",
					HelperText:  codeError,
					HelperType:  codeHelperType,
					InputProps: []g.Node{
						h.Value(codeValue),
						h.AutoComplete("off"),
					},
				}),
			),

			g.If(
				p.ErrorText != "",
				h.Div(
					h.Class("error-msg"),
					g.Text(p.ErrorText),
				),
			),

			g.If(
				p.SuccessText != "",
				h.Div(
					h.Class("success-msg"),
					g.Text(p.SuccessText),
				),
			),

			components.Button(
				&components.ButtonProps{},
				h.Type("submit"),
				g.Text("Add SKU Item"),
			),
		),
	})

	return layout.Page(layout.PageProps{
		Ctx:   p.Ctx,
		Title: "Add New SKU Configuration",
		Breadcrumbs: []layout.Breadcrumb{
			layout.HomeBreadcrumb,
			{
				IconIdentifier: "package-variant-closed",
				Title:          "Stock Items",
				URLPart:        "stock-items",
			},
			{
				IconIdentifier: "barcode-scan",
				Title:          "SKU Configuration",
				URLPart:        "sku-config",
			},
			{
				IconIdentifier: "plus",
				Title:          "Add",
			},
		},
		Content: content,
		AppendHead: []g.Node{
			components.InlineStyle("/internal/views/userview/add_user_page.css"),
			components.InlineStyle("/internal/views/stockitemview/add_sku_item.css"),
		},
		AppendBody: []g.Node{
			components.InlineScript("/internal/views/stockitemview/add_stock_item.js"),
		},
	})
}

var SKUFields = []components.SearchSelectOption{
	{Label: "Product Type", Value: "ProductType"},
	{Label: "Yarn Type", Value: "YarnType"},
	{Label: "Style Number", Value: "StyleNumber"},
	{Label: "Colour", Value: "Colour"},
	{Label: "Toe Closing", Value: "ToeClosing"},
	{Label: "Size", Value: "Size"},
}
