package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
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
		classes,
	)
}
