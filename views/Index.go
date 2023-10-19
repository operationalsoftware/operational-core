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
		// {
		// 	Text:     "Test",
		// 	UrlToken: "test",
		// },
	}

	indexContent := g.Group([]g.Node{
		h.H1(g.Text("Operational Core Home")),
		// Buttons
		h.Div(
			h.Class("button-group"),
			g.Raw(`
			<style>
			  me {
					display: flex;
					justify-content: center;
					align-items: center;
					flex-wrap: wrap;
					gap: 1rem;
				}
			</style>`,
			),
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

			h.Div(
				h.Class("loaders"),
				g.Raw(`
				<style>
			  me {
					display: flex;
					justify-content: center;
					align-items: center;
					flex-wrap: wrap;
					gap: 1rem;
					margin-top: 1rem;
				}
				</style>`,
				),
				o.LoadingSpinner(o.LoadingSpinnerXs),
				o.LoadingSpinner(o.LoadingSpinnerSm),
				o.LoadingSpinner(""),
				o.LoadingSpinner(o.LoadingSpinnerLg),
			),
		),
	})

	return layout.Page(layout.PageParams{
		Title:   "Home",
		Content: indexContent,
		Crumbs:  crumbs,
	})
}
