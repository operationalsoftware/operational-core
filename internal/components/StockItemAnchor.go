package components

import (
	"net/url"

	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"
)

func StockItemAnchor(stockCode string, children ...g.Node) g.Node {

	hrefStr := "/stock/" + url.PathEscape(stockCode)

	return h.A(
		h.Href(hrefStr),
		g.Group(children),
		g.Text(stockCode),
	)
}
