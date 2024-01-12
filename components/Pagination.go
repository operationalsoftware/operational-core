package components

import (
	"fmt"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type PaginationProps struct {
	CurrentPage int
	TotalPages  int
	PageSize    int
}

func Pagination(p *PaginationProps) g.Node {
	classes := c.Classes{
		"pagination": true,
	}

	if p.CurrentPage > p.TotalPages {
		p.CurrentPage = p.TotalPages
	}

	// If there is only one page, don't render pagination.
	if p.TotalPages <= 1 {
		classes["pagination--hidden"] = true
	}

	pageArray := make([]int, p.TotalPages)

	for i := range pageArray {
		pageArray[i] = i
	}

	return h.Div(
		classes,
		h.DataAttr("current-page", fmt.Sprintf("%d", p.CurrentPage)),
		h.DataAttr("total-pages", fmt.Sprintf("%d", p.TotalPages)),
		h.Ul(
			h.Class("pagination__list"),
			// Chevron left
			h.Li(
				h.Class("pagination__btn pagination__btn--left"),
				h.A(
					Icon(&IconProps{
						Identifier: "chevron-left",
						Classes: c.Classes{
							"pagination__icon": true,
						},
					}),
				),
			),
			g.Group(g.Map(pageArray, func(i int) g.Node {
				page := i + 1
				classes := c.Classes{
					"pagination__item": true,
				}
				if page == p.CurrentPage {
					classes["pagination__item--current"] = true
				}
				return h.Li(
					classes,
					h.A(
						g.Text(fmt.Sprintf("%d", page)),
					),
				)
			})),
			// Chevron right
			h.Li(
				h.Class("pagination__btn pagination__btn--right"),
				h.A(
					Icon(&IconProps{
						Identifier: "chevron-right",
						Classes: c.Classes{
							"pagination__icon": true,
						},
					}),
				),
			),
		),

		InlineStyle(Assets, "/Pagination.css"),
		InlineScript(Assets, "/Pagination.js"),
	)

}
