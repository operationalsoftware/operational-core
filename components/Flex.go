package components

import (
	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type Direction string

const (
	Row    Direction = "row"
	Column Direction = "column"
)

type Justify string

const (
	JustifyStart        Justify = "flex-start"
	JustifyEnd          Justify = "flex-end"
	JustifyCenter       Justify = "center"
	JustifySpaceBetween Justify = "space-between"
	JustifySpaceAround  Justify = "space-around"
	JustifySpaceEvenly  Justify = "space-evenly"
	Inherit             Justify = "inherit"
	Initial             Justify = "initial"
)

type Align string

const (
	AlignStart    Align = "flex-start"
	AlignEnd      Align = "flex-end"
	AlignCenter   Align = "center"
	AlignBaseline Align = "baseline"
	AlignStretch  Align = "stretch"
)

type FlexProps struct {
	wrap      bool
	direction Direction
	justify   Justify
	align     Align
}

func Flex(children ...g.Node) g.Node {
	return h.Div(
		h.Class("flex"),
		g.Group(children),
	)
}
