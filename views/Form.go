package views

import (
	o "operationalcore/components"
	"operationalcore/layout"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

var formCrumb layout.Crumb = layout.Crumb{
	Text:     "Form",
	UrlToken: "form",
}

func Form() g.Node {
	crumbs := []layout.Crumb{
		formCrumb,
	}

	formContent := g.Group([]g.Node{
		h.H1(g.Text("Form Page")),
		o.Form(
			o.Input(&o.InputProps{
				Label:       "Text",
				Name:        "text",
				Placeholder: "Enter text",
			}),

			o.Select(&o.SelectProps{
				Name: "single-select",
				Options: []o.Option{
					{Value: "1", Label: "One"},
					{Value: "2", Label: "Two"},
					{Value: "3", Label: "Three"},
					{Value: "hello-world", Label: "Hello world"},
				},
				Multiple: false,
			}),

			o.Select(&o.SelectProps{
				Name: "multi-select",
				Options: []o.Option{
					{Value: "1", Label: "One"},
					{Value: "2", Label: "Two"},
					{Value: "3", Label: "Three"},
					{Value: "hello-world", Label: "Hello world"},
				},
				Multiple: true,
			}),

			o.SearchSelect(&o.SearchSelectProps{
				Name: "search-select",
				Options: []o.Option{
					{Value: "1", Label: "One"},
					{Value: "2", Label: "Two"},
					{Value: "3", Label: "Three"},
					{Value: "hello-world", Label: "Hello world"},
				},
				OptionUrl: "/options",
				Multiple:  true,
			}),

			o.Button(&o.ButtonProps{}, h.Type("submit"), g.Text("Submit")),
		),
		o.InlineScript(Assets, "/Form.js"),
	})

	return layout.Page(layout.PageParams{
		Title:   "Form",
		Content: formContent,
		Crumbs:  crumbs,
	})
}
