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

// Deprecated: Use CSS size class strings directly on spinner markup.
type LoadingSpinnerSize string

// Deprecated: Build spinner markup directly and apply "loading-spinner" plus size classes.
type LoadingSpinnerProps struct {
	Size LoadingSpinnerSize
}

// Deprecated: Use h.Div directly and apply "loading-spinner" and size classes instead of this wrapper component.
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
