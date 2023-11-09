package views

import (
	"fmt"
	o "operationalcore/components"
	"operationalcore/layout"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type CustomDataRow struct {
	col1 string
	col2 string
	col3 string
	col4 int
}

func (t CustomDataRow) Render() map[string]o.RenderedCell {
	return map[string]o.RenderedCell{
		"col-1": {
			Content: g.Text(t.col1),
			Classes: c.Classes{
				"table-cell": true,
			},
		},
		"col-2": {
			Content: g.Text(t.col2),
			Classes: c.Classes{
				"table-cell": true,
			},
		},
		"col-3": {
			Content: g.Text(t.col3),
			Classes: c.Classes{
				"table-cell": true,
			},
		},
		"col-4": {
			Content: g.Text(fmt.Sprint(t.col4)),
			Classes: c.Classes{
				"table-cell": true,
			},
		},
	}
}

var columns = []o.TableColumn{
	{
		Name: "Column 1",
		Key:  "col-1",
	},
	{
		Name: "Column 2",
		Key:  "col-2",
	},
	{
		Name: "Column 3",
		Key:  "col-3",
	},
	{
		Name: "Column 4",
		Key:  "col-4",
	},
}

var props = &o.TableProps{
	Columns: columns,
	Data:    data,
}

var data = []o.TableRowRenderer{
	CustomDataRow{
		col1: "Data 1",
		col2: "Data 2",
		col3: "Data 3",
		col4: 4,
	},
	CustomDataRow{
		col1: "Data 1",
		col2: "Data 2",
		col3: "Data 3",
		col4: 8,
	},
}

var indexCrumb layout.Crumb = layout.Crumb{
	Text:     "Home",
	UrlToken: "",
}

func Index() g.Node {
	crumbs := []layout.Crumb{
		indexCrumb,
	}

	indexContent := g.Group([]g.Node{
		h.H1(g.Text("Operational Core Home")),
		// Buttons
		o.Card(
			g.Raw(`<style>
			me {
				display: flex;
				justify-content: center;
				align-items: center;
				gap: 1rem;
			}
		</style>`),
			o.Button(&o.ButtonProps{
				ButtonType: o.ButtonPrimary,
				Size:       o.ButtonLg,
				Loading:    false,
				Disabled:   false,
			},
				g.Text("Primary"),
			),
			o.Button(&o.ButtonProps{
				ButtonType: o.ButtonSecondary,
				Size:       "",
				Loading:    false,
				Disabled:   false,
			},
				g.Text("Secondary"),
			),
			o.Button(&o.ButtonProps{
				ButtonType: o.ButtonWarning,
				Size:       "",
				Loading:    false,
				Disabled:   false,
			},
				g.Text("Warning"),
			),
			o.Button(&o.ButtonProps{
				ButtonType: o.ButtonSuccess,
				Size:       o.ButtonSm,
				Loading:    false,
				Disabled:   false,
			},
				g.Text("Success"),
			),
			o.Button(&o.ButtonProps{
				ButtonType: o.ButtonSecondary,
				Size:       o.ButtonSm,
				Loading:    true,
				Disabled:   false,
			},
				g.Text("Loading"),
			),
		),
		// Loaders
		o.Card(
			g.Raw(`<style>
			  me {
					display: flex;
					justify-content: center;
					align-items: center;
					gap: 1rem;
				}
			</style>`),
			o.LoadingSpinner(o.LoadingSpinnerSm),
			o.LoadingSpinner(""),
			o.LoadingSpinner(o.LoadingSpinnerLg),
			o.LoadingSpinner(o.LoadingSpinnerXl),
		),
		// Tooltips
		o.Card(
			g.Raw(`<style>
			  me {
					display: flex;
					justify-content: flex-start;
					align-items: center;
				}
			</style>`),
			o.Tooltip(&o.TooltipProps{
				Text:     "Tooltip Top Content",
				Position: "",
			},
				o.Button(&o.ButtonProps{
					ButtonType: o.ButtonSecondary,
					Size:       o.ButtonSm,
					Loading:    false,
					Disabled:   false,
				},
					g.Text("Tooltip Top"),
				),
			),
		),
		// Statistic
		o.Card(
			g.Raw(`<style>
			  me {
					display: flex;
					justify-content: space-between;
					align-items: center;
				}
			</style>`),
			o.Statistic(&o.StatisticProps{
				Heading: "Active users",
				Value:   "1,234",
			}),
			o.Statistic(&o.StatisticProps{
				Heading: "Bank Balance",
				Value:   "1,324324",
			}),
			o.Statistic(&o.StatisticProps{
				Heading: "Total Withdraws",
				Value:   "50",
			}),
		),
		// Modal
		o.Button(
			&o.ButtonProps{
				ButtonType: o.ButtonPrimary,
				Size:       o.ButtonSm,
				Loading:    false,
				Disabled:   false,
			},
			g.Group([]g.Node{
				g.Text("Open Modal"),
				h.ID("open-modal"),
			}),
		),
		o.Modal(&o.ModalProps{
			Title:         "Sample Modal",
			FooterContent: h.P(g.Text("This is modal footer")),
		},
			h.P(g.Text("This is modal content")),
		),
		// Progress
		o.Card(
			o.Progress(&o.ProgressProps{
				Percentage: 40,
			}, h.ID("progress-1")),
			o.Progress(&o.ProgressProps{
				Percentage: 50,
				Type:       o.ProgressTypeWarning,
			}, h.ID("progress-2")),
			o.Progress(&o.ProgressProps{
				Percentage: 60,
				Type:       o.ProgressTypeDanger,
			}, h.ID("progress-3")),
		),
		// Inputs
		o.Card(
			o.Form(
				o.Input(&o.InputProps{
					Size:        "",
					Placeholder: "Small",
					Name:        "small",
					Label:       "Small Input",
				}),
				o.InputNumber(&o.InputNumberProps{
					Size:        "",
					Placeholder: "Small Input Number",
					Label:       "Small Input Number",
					Name:        "small-number",
					HelperText:  "This is helper text",
					HelperType:  "",
				}),
				o.Textarea(&o.TextareaProps{
					Name:        "textarea",
					Label:       "Textarea",
					Placeholder: "Write anything",
				}),
				// Radio
				o.Radio(&o.RadioProps{
					Name:  "radio-1",
					Label: "Radio 1",
				}),
				o.Radio(&o.RadioProps{
					Name:  "radio-2",
					Label: "Radio 2",
				}),
				o.Radio(&o.RadioProps{
					Name:  "radio-3",
					Label: "Radio 3",
				}),
				// Checkbox
				o.Checkbox(
					&o.CheckboxProps{
						Name:    "checkbox-1",
						Label:   "Checkbox 1",
						Value:   "1",
						Checked: true,
					},
				),
			),
		),
		// Pop confirms
		o.Card(
			g.Raw(`<style>
			  me {
					display: flex;
					justify-content: flex-start;
					align-items: center;
				}
			</style>`),
			o.Popconfirm(&o.PopconfirmProps{
				Id:      "popconfirm-1",
				Icon:    "account",
				Heading: "Hello world",
				Text:    "Are you sure you want to do this!",
				Yes:     "Yes",
				No:      "No",
			},
				o.Button(
					&o.ButtonProps{
						ButtonType: o.ButtonPrimary,
						Size:       o.ButtonSm,
						Loading:    false,
						Disabled:   false,
						Classes: c.Classes{
							"popconfirm-trigger": true,
						},
					},
					g.Text("Open Popconfirm"),
				),
			),
		),
		// Upload button
		o.Card(
			o.UploadButton(
				&o.UploadButtonProps{
					Id: "upload-button-1",
				},
			),
		),
		// Table
		o.Card(
			o.Table(props),
		),
		// Slider
		o.Card(
			o.Slider(&o.SliderProps{
				Min: "0",
				Max: "10",
			}),
		),
		o.InlineScript(Assets, "/Index.js"),
	})

	return layout.Page(layout.PageParams{
		Title:   "Home",
		Content: indexContent,
		Crumbs:  crumbs,
	})
}
