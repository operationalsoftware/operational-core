package views

import (
	o "operationalcore/components"
	"operationalcore/layout"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

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
			g.Text("Open Modal"),
			h.ID("open-modal"),
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
				Percentage: 0,
			}, h.ID("progress-1")),
			o.Progress(&o.ProgressProps{
				Percentage: 0,
				Type:       o.ProgressTypeWarning,
			}, h.ID("progress-2")),
			o.Progress(&o.ProgressProps{
				Percentage: 0,
				Type:       o.ProgressTypeDanger,
			}, h.ID("progress-3")),
		),
		// Inputs
		o.Card(
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
		),
	})

	return layout.Page(layout.PageParams{
		Title:   "Home",
		Content: indexContent,
		Crumbs:  crumbs,
	})
}
