package components

import (
	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	h "maragu.dev/gomponents/html"
)

const (
	LoadingSpinnerSm LoadingSpinnerSize = "sm"
	LoadingSpinnerMd LoadingSpinnerSize = "md"
	LoadingSpinnerLg LoadingSpinnerSize = "lg"
	LoadingSpinnerXl LoadingSpinnerSize = "xl"
)

type LoadingSpinnerSize string

type LoadingSpinnerProps struct {
	Size LoadingSpinnerSize
}

func LoadingSpinner(p *LoadingSpinnerProps) g.Node {
	classes := c.Classes{
		"loading-spinner": true,
	}

	if p.Size == "" {
		p.Size = LoadingSpinnerMd
	}

	classes[string(p.Size)] = true
	return h.Div(
		classes,
	)
}
