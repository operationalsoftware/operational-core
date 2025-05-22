package components

import (
	"net/url"

	g "github.com/maragudk/gomponents"
	h "github.com/maragudk/gomponents/html"
)

func StockItemAnchor(stockCode string, children ...g.Node) g.Node {

	hrefStr := "/stock/" + url.PathEscape(stockCode)

	return h.A(
		h.Href(hrefStr),
		g.Group(children),
		g.Text(stockCode),
	)
}
