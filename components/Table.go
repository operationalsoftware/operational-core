package components

import (
	"app/utils"
	"fmt"
	"net/url"

	g "github.com/maragudk/gomponents"
	hx "github.com/maragudk/gomponents-htmx"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type RenderedCell struct {
	Content    g.Node
	Classes    c.Classes
	Attributes []g.Node
}

type TableRowRenderer interface {
	Render() map[string]RenderedCell
}

type TableColumn struct {
	Name string
	Key  string
}

type TableColumns []TableColumn

func getSortIconIdentifier(sort utils.Sort, key string) string {
	sortDirection := sort.GetSortDirection(key)
	if sortDirection == "" {
		return "arrow-up-down"
	} else if sortDirection == utils.SortDirectionAsc {
		return "arrow-up"
	}
	return "arrow-down"
}

func generateSortString(currentSort utils.Sort, key string) string {

	// create a copy of the current sort
	// this is important as the current sort may be passed to this function several
	// times with different keys, and we don't want to mutate the original sort
	sort := make(utils.Sort, len(currentSort))
	copy(sort, currentSort)

	// If the key is not in the sort, add it with the default sort direction (asc)
	// If the key is in the sort:
	//   - if it's the last key in the sort and is sorted by desc, remove it
	//   - else, reverse the sort direction

	isInSort := false
	for i, s := range sort {
		if s.Key == key {
			isInSort = true
			if i == len(sort)-1 && s.Sort == utils.SortDirectionDesc {
				// remove it
				sort = append(sort[:i], sort[i+1:]...)
			} else {
				// reverse the sort direction
				if s.Sort == utils.SortDirectionAsc {
					sort[i].Sort = utils.SortDirectionDesc
				} else {
					sort[i].Sort = utils.SortDirectionAsc
				}
			}
		}
	}

	if !isInSort {
		sort = append(sort, utils.SortItem{
			Key:  key,
			Sort: utils.SortDirectionAsc,
		})
	}

	return sort.EncodeQueryParam()
}

func generateSortButtonURL(
	path string,
	urlQuery url.Values,
	sortQueryKey string,
	sort utils.Sort,
	columnKey string,
) string {
	// make a copy of the url urlQuery
	urlQueryCopy := url.Values{}
	for k, v := range urlQuery {
		urlQueryCopy[k] = v
	}

	// generate the sort string
	sortString := generateSortString(sort, columnKey)

	if sortString == "" {
		// remove the sort key from the urlQuery
		urlQueryCopy.Del(sortQueryKey)
	} else {
		// set the sort string in the urlQuery
		urlQueryCopy.Set(sortQueryKey, sortString)
	}

	// return the url
	encodedQuery := urlQueryCopy.Encode()
	if encodedQuery != "" {
		encodedQuery = "?" + encodedQuery
	}
	return path + encodedQuery
}

type sortButtonProps struct {
	hxGetPath    string
	hxSelect     string
	urlQuery     url.Values
	sortQueryKey string
	sort         utils.Sort
	columnKey    string
}

func sortButton(p *sortButtonProps) g.Node {

	iconIdentifier := getSortIconIdentifier(p.sort, p.columnKey)
	hxGetURL := generateSortButtonURL(
		p.hxGetPath,
		p.urlQuery,
		p.sortQueryKey,
		p.sort, p.columnKey,
	)

	sortPosition := p.sort.GetSortPosition(p.columnKey)

	return h.Button(
		h.Class("table-sort-button"),
		hx.Get(hxGetURL),
		hx.Target("closest .table-container"),
		hx.PushURL("true"),
		g.If(p.hxSelect != "", hx.Select(p.hxSelect)),
		Icon(&IconProps{
			Identifier: iconIdentifier,
		}),
		g.If(sortPosition != -1, h.Span(g.Text(fmt.Sprintf("%d", sortPosition+1)))),
	)
}

type tableHeadProps struct {
	columns      TableColumns
	sortableKeys []string
	hxGetPath    string
	hxSelect     string
	urlQuery     url.Values
	sortQueryKey string
}

func tableHead(p *tableHeadProps) g.Node {
	sort := utils.Sort{}
	sort.ParseQueryParam(p.urlQuery.Get(p.sortQueryKey), p.sortableKeys)

	headerCells := g.Group(g.Map(p.columns, func(c TableColumn) g.Node {

		var sortButtonNode g.Node
		sortable := false
		for _, k := range p.sortableKeys {
			if k == c.Key {
				sortable = true
				break
			}
		}
		if sortable {
			sortButtonNode = sortButton(&sortButtonProps{
				hxGetPath:    p.hxGetPath,
				hxSelect:     p.hxSelect,
				urlQuery:     p.urlQuery,
				sortQueryKey: p.sortQueryKey,
				sort:         sort,
				columnKey:    c.Key,
			})
		}

		return h.Th(
			h.Class("table-head"),
			h.Span(g.Text(c.Name)),
			g.If(sortable, sortButtonNode),
		)
	}))

	return h.THead(h.Tr(headerCells))
}

func renderRows(p *TableProps) g.Node {
	return g.Group(g.Map(p.Data, func(d TableRowRenderer) g.Node {
		return h.Tr(
			h.Class("table-row"),
			g.Group(g.Map(p.Columns, func(col TableColumn) g.Node {
				renderKey := col.Key
				if renderKey == "" {
					panic("TableColumn.Key must be set")
				}

				renderedCell := d.Render()[renderKey]

				if renderedCell.Classes == nil {
					renderedCell.Classes = c.Classes{}
				}

				renderedCell.Classes["table-cell"] = true

				return h.Td(
					renderedCell.Classes,
					g.Group(renderedCell.Attributes),
					renderedCell.Content,
				)
			})),
		)
	}))
}

type TableProps struct {
	Columns      TableColumns
	SortableKeys []string
	Data         []TableRowRenderer
	UrlQuery     url.Values
	SortQueryKey string
	HXGetPath    string
	HXSelect     string
}

func Table(p *TableProps, children ...g.Node) g.Node {
	classes := c.Classes{
		"table-container": true,
	}

	if p.SortQueryKey == "" {
		p.SortQueryKey = "sort"
	}

	if p.HXGetPath == "" {
		panic("TableProps.HXGetPath must be set")
	}
	if p.HXSelect == "" {
		panic("TableProps.HXSelect must be set")
	}

	return h.Div(
		classes,
		h.Table(
			tableHead(&tableHeadProps{
				columns:      p.Columns,
				sortableKeys: p.SortableKeys,
				hxGetPath:    p.HXGetPath,
				hxSelect:     p.HXSelect,
				urlQuery:     p.UrlQuery,
				sortQueryKey: p.SortQueryKey,
			}),

			h.TBody(
				renderRows(p),
			),
		),
		g.Group(children),
		InlineStyle("/components/Table.css"),
	)
}
