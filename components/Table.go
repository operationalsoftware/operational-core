package components

import (
	"app/utils"
	"fmt"
	"math"

	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	h "github.com/maragudk/gomponents/html"
)

type TableCell struct {
	Contents   g.Node
	Classes    c.Classes
	Attributes []g.Node
}

type TableRow struct {
	Cells      []TableCell
	Classes    c.Classes
	Attributes []g.Node
	SubRows    []TableRow
}

type TableRows []TableRow

type TableColumn struct {
	TitleContents g.Node
	Classes       c.Classes
	Attributes    []g.Node
	SortKey       string // unset if column isn't sortable
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

type sortRadioProps struct {
	onChange      string
	sortQueryKey  string
	sort          utils.Sort
	columnSortKey string
}

func sortRadio(p *sortRadioProps) g.Node {

	iconIdentifier := getSortIconIdentifier(p.sort, p.columnSortKey)
	sortPosition := p.sort.GetSortPosition(p.columnSortKey)

	return h.Label(
		h.Input(
			h.Type("radio"),
			h.Name(p.sortQueryKey),
			h.Value(generateSortString(p.sort, p.columnSortKey)),
			g.Attr("onchange", p.onChange),
		),
		Icon(&IconProps{
			Identifier: iconIdentifier,
		}),
		g.If(sortPosition != -1, h.Span(g.Text(fmt.Sprintf("%d", sortPosition+1)))),
	)
}

type tableHeadProps struct {
	columns      TableColumns
	sortableKeys []string
	sort         utils.Sort
	sortQueryKey string
	onChange     string
}

func tableHead(p *tableHeadProps) g.Node {
	headerCells := g.Group(g.Map(p.columns, func(col TableColumn) g.Node {

		var sortRadioNode g.Node
		sortable := col.SortKey != ""
		if sortable {
			sortRadioNode = sortRadio(&sortRadioProps{
				onChange:      p.onChange,
				sortQueryKey:  p.sortQueryKey,
				sort:          p.sort,
				columnSortKey: col.SortKey,
			})
		}

		if col.Classes == nil {
			col.Classes = c.Classes{}
		}

		col.Classes["table-head"] = true

		return h.Th(
			col.Classes,
			g.Group(col.Attributes),
			h.Span(col.TitleContents),
			g.If(sortable, sortRadioNode),
		)
	}))

	return h.THead(h.Tr(headerCells))
}

func renderRows(rows TableRows) g.Node {
	return g.Group(g.Map(rows, func(row TableRow) g.Node {

		if row.Classes == nil {
			row.Classes = c.Classes{}
		}

		row.Classes["table-row"] = true

		return h.Tr(
			row.Classes,
			g.Group(row.Attributes),
			g.Group(g.Map(row.Cells, func(cell TableCell) g.Node {
				if cell.Classes == nil {
					cell.Classes = c.Classes{}
				}

				return h.Td(
					cell.Classes,
					g.Group(cell.Attributes),
					cell.Contents,
				)
			})),
			g.If(
				row.SubRows != nil,
				renderRows(row.SubRows),
			),
		)
	}))
}

type TablePaginationProps struct {
	TotalRecords        int
	CurrentPage         int
	CurrentPageQueryKey string
	PageSize            int
	PageSizeQueryKey    string
}

func TablePagination(p TablePaginationProps, onChange string) g.Node {

	if p.PageSizeQueryKey == "" {
		p.PageSizeQueryKey = "PageSize"
	}
	if p.CurrentPageQueryKey == "" {
		p.CurrentPageQueryKey = "Page"
	}

	// calculate the total pages
	totalPages := int(math.Ceil(float64(p.TotalRecords) / float64(p.PageSize)))

	// if the current page is less than 1, set it to 1
	if p.CurrentPage < 1 {
		p.CurrentPage = 1
	}
	// if the current page is greater than the total pages, set it to the total pages
	if p.CurrentPage > totalPages {
		p.CurrentPage = totalPages
	}

	// define a pages array
	currentMinus10 := p.CurrentPage - 10
	if currentMinus10 < 1 {
		currentMinus10 = 0
	}
	currentPlus10 := p.CurrentPage + 10
	if currentPlus10 > totalPages {
		currentPlus10 = 0
	}
	pages := []int{}
	// add pages -5 to +5 from the current page
	for i := p.CurrentPage - 5; i <= p.CurrentPage+5; i++ {
		if i > 0 && i <= totalPages {
			pages = append(pages, i)
		}
	}

	pageRadio := func(page int, label string, disabled bool) g.Node {
		classes := c.Classes{
			"table-pagination-page": true,
		}
		if disabled {
			classes["disabled"] = true
		}
		if page == p.CurrentPage && label == "" {
			classes["active"] = true
		}

		return h.Label(
			classes,
			h.Input(
				h.Type("radio"),
				h.Name(p.CurrentPageQueryKey),
				g.If(disabled, h.Disabled()),
				h.Value(fmt.Sprintf("%d", page)),
				h.StyleAttr("display: none"),
				g.Attr("onchange", onChange),
			),
			g.If(label != "", g.Text(label)),
			g.If(label == "", g.Text(fmt.Sprintf("%d", page))),
		)
	}

	return h.Div(
		h.Class("table-pagination"),

		// use radio buttons for changing the page
		h.Div(
			h.Class("table-pagination-pages"),
			// start
			pageRadio(1, "\u21E4", p.CurrentPage == 1),
			// skip 10 previous
			g.If(currentMinus10 > 0, pageRadio(currentMinus10, "\u226A", false)),
			// previous
			pageRadio(p.CurrentPage-1, "\u003C", p.CurrentPage == 1),
			g.Group(g.Map(pages, func(i int) g.Node {
				return pageRadio(i, "", false)
			})),
			// next
			pageRadio(p.CurrentPage+1, "\u003E", p.CurrentPage == totalPages),
			// skip 10 next
			g.If(currentPlus10 > 0, pageRadio(currentPlus10, "\u226B", false)),
			// end
			pageRadio(totalPages, "\u21E5", p.CurrentPage == totalPages),
		),
	)
}

type TableProps struct {
	Classes      c.Classes
	Columns      TableColumns
	SortableKeys []string
	Rows         TableRows
	Sort         utils.Sort
	SortQueryKey string
	Pagination   TablePaginationProps
	OnChange     string
}

func Table(p *TableProps, children ...g.Node) g.Node {
	if p.Classes == nil {
		p.Classes = c.Classes{}
	}
	p.Classes["table"] = true

	if p.SortQueryKey == "" {
		p.SortQueryKey = "Sort"
	}

	return h.Div(
		p.Classes,
		g.Group(children),
		g.If(
			p.Pagination != (TablePaginationProps{}),
			TablePagination(p.Pagination, p.OnChange),
		),
		h.Div(
			h.Class("table-scroll"),
			h.Table(
				tableHead(&tableHeadProps{
					columns:      p.Columns,
					onChange:     p.OnChange,
					sortableKeys: p.SortableKeys,
					sort:         p.Sort,
					sortQueryKey: p.SortQueryKey,
				}),

				h.TBody(
					renderRows(p.Rows),
				),
			),
			// here we add a hidden radio inputs which serve the purpose of
			// preserving sort, page and page size state in the URL
			g.If(
				p.Sort != nil,
				h.Input(
					h.Type("radio"),
					h.Checked(),
					h.Name(p.SortQueryKey),
					h.Value(p.Sort.EncodeQueryParam()),
					h.StyleAttr("display: none"),
				),
			),
			g.If(
				p.Pagination != (TablePaginationProps{}),
				g.Group([]g.Node{
					h.Input(
						h.Type("radio"),
						h.Checked(),
						h.Name(p.Pagination.CurrentPageQueryKey),
						h.Value(fmt.Sprintf("%d", p.Pagination.CurrentPage)),
						h.StyleAttr("display: none"),
					),
					h.Input(
						h.Type("radio"),
						h.Checked(),
						h.Name(p.Pagination.PageSizeQueryKey),
						h.Value(fmt.Sprintf("%d", p.Pagination.PageSize)),
						h.StyleAttr("display: none"),
					),
				}),
			),
		),
		g.If(
			p.Pagination != (TablePaginationProps{}),
			TablePagination(p.Pagination, p.OnChange),
		),
	)
}
