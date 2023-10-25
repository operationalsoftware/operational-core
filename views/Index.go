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
				Text:       "Button LG",
				ButtonType: o.ButtonPrimary,
				Size:       o.ButtonLg,
				Loading:    false,
				Disabled:   false,
			},
			),
			o.Button(&o.ButtonProps{
				Text:       "Button MD",
				ButtonType: o.ButtonSecondary,
				Size:       "",
				Loading:    false,
				Disabled:   false,
			},
			),
			o.Button(&o.ButtonProps{
				Text:       "Button Warning",
				ButtonType: o.ButtonWarning,
				Size:       "",
				Loading:    false,
				Disabled:   false,
			},
			),
			o.Button(&o.ButtonProps{
				Text:       "Button SM",
				ButtonType: o.ButtonSuccess,
				Size:       o.ButtonSm,
				Loading:    false,
				Disabled:   false,
			},
			),
			o.Button(&o.ButtonProps{
				Text:       "Button Loading Data Attribute",
				ButtonType: o.ButtonSecondary,
				Size:       o.ButtonSm,
				Loading:    true,
				Disabled:   false,
			},
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
				Text:     "Tooltip Right Content",
				Position: "",
			},
				o.Button(&o.ButtonProps{
					Text:       "Trigger Tooltip",
					ButtonType: o.ButtonSecondary,
					Size:       o.ButtonSm,
					Loading:    false,
					Disabled:   false,
				},
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
				Text:       "Trigger Modal",
				ButtonType: o.ButtonPrimary,
				Size:       o.ButtonSm,
				Loading:    false,
				Disabled:   false,
				Id:         "open-modal",
			},
		),
		o.Modal(&o.ModalProps{
			Title:         "Sample Modal",
			FooterContent: h.P(g.Text("This is modal footer")),
		},
			h.P(g.Text("This is modal content")),
		),
	})

	return layout.Page(layout.PageParams{
		Title:   "Home",
		Content: indexContent,
		Crumbs:  crumbs,
	})
}
