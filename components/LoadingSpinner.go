package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type LoadingSpinnerSize string

const (
	LoadingSpinnerXs LoadingSpinnerSize = "xs"
	LoadingSpinnerSm LoadingSpinnerSize = "sm"
	LoadingSpinnerLg LoadingSpinnerSize = "lg"
)

func LoadingSpinner(size LoadingSpinnerSize) g.Node {
	classes := c.Classes{
		"loading-spinner": true,
	}
	if size != "" {
		classes[string(size)] = true
	}
	return h.Div(
		InlineStyle(Assets, "/LoadingSpinner.css"),
		classes,
	)
}
