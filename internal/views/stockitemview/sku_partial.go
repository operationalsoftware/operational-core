package stockitemview

import (
	"app/internal/components"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

type SKUPartialProps struct {
	SKU string
}

func StockCodePartial(data SKUPartialProps) g.Node {
	content := components.Input(&components.InputProps{
		Label:       "Stock Code (SKU)",
		Name:        "StockCode",
		Placeholder: "Stock code",
		InputProps: []g.Node{
			h.Value(data.SKU),
			h.AutoComplete("off"),
			h.ReadOnly(),
		},
	})

	return content
}
