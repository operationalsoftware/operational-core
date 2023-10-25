package components

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type LoadingSpinnerSize string

const (
	LoadingSpinnerSm LoadingSpinnerSize = "sm"
	LoadingSpinnerMd LoadingSpinnerSize = "md"
	LoadingSpinnerLg LoadingSpinnerSize = "lg"
	LoadingSpinnerXl LoadingSpinnerSize = "xl"
)

func LoadingSpinner(size LoadingSpinnerSize) g.Node {
	classes := c.Classes{
		"loading-spinner": true,
	}

	if size == "" {
		size = LoadingSpinnerMd
	}

	classes[string(size)] = true
	return h.Div(
		InlineStyle(Assets, "/LoadingSpinner.css"),
		classes,
	)
}
