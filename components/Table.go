package components

import (
	"app/internal/appsort"
	"fmt"
	"math"
	"strconv"

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

func getSortIconIdentifier(sort appsort.Sort, key string) string {
	sortDirection := sort.GetDirection(key)
	if sortDirection == "" {
		return "arrow-up-down"
	} else if sortDirection == appsort.DirectionAsc {
		return "arrow-up"
	}
	return "arrow-down"
}

func generateSortString(currentSort appsort.Sort, key string) string {

	// create a copy of the current sort
	// this is important as the current sort may be passed to this function several
	// times with different keys, and we don't want to mutate the original sort
	sort := make(appsort.Sort, len(currentSort))
	copy(sort, currentSort)

	// If the key is not in the sort, add it with the default sort direction (asc)
	// If the key is in the sort:
	//   - if it's the last key in the sort and is sorted by desc, remove it
	//   - else, reverse the sort direction

	isInSort := false
	for i, s := range sort {
		if s.Key == key {
			isInSort = true
			if i == len(sort)-1 && s.Sort == appsort.DirectionDesc {
				// remove it
				sort = append(sort[:i], sort[i+1:]...)
			} else {
				// reverse the sort direction
				if s.Sort == appsort.DirectionAsc {
					sort[i].Sort = appsort.DirectionDesc
				} else {
					sort[i].Sort = appsort.DirectionAsc
				}
			}
		}
	}

	if !isInSort {
		sort = append(sort, appsort.SortItem{
			Key:  key,
			Sort: appsort.DirectionAsc,
		})
	}

	return sort.EncodeQueryParam()
}

type sortRadioProps struct {
	sortQueryKey  string
	sort          appsort.Sort
	columnSortKey string
}

func sortRadio(p *sortRadioProps) g.Node {

	iconIdentifier := getSortIconIdentifier(p.sort, p.columnSortKey)
	sortIndex := p.sort.GetIndex(p.columnSortKey)

	return h.Label(
		h.Input(
			h.Type("radio"),
			h.Name(p.sortQueryKey),
			h.Value(generateSortString(p.sort, p.columnSortKey)),
			g.Attr("onchange", "submitTableForm(this.form)"),
		),
		Icon(&IconProps{
			Identifier: iconIdentifier,
		}),
		g.If(sortIndex != -1, h.Span(g.Text(fmt.Sprintf("%d", sortIndex+1)))),
	)
}

type tableHeadProps struct {
	columns      TableColumns
	sort         appsort.Sort
	sortQueryKey string
}

func tableHead(p *tableHeadProps) g.Node {
	headerCells := g.Group(g.Map(p.columns, func(col TableColumn) g.Node {

		var sortRadioNode g.Node
		sortable := col.SortKey != ""
		if sortable {
			sortRadioNode = sortRadio(&sortRadioProps{
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
	PageSizeOptions     []int
	PageSize            int
	PageSizeQueryKey    string
}

func TablePagination(p *TablePaginationProps) g.Node {

	if p == nil {
		// has been called with nil pointer
		return nil
	}

	if p.PageSizeQueryKey == "" {
		p.PageSizeQueryKey = "PageSize"
	}
	if p.CurrentPageQueryKey == "" {
		p.CurrentPageQueryKey = "Page"
	}
	if p.PageSizeOptions == nil {
		p.PageSizeOptions = []int{25, 50, 100, 200, 500}
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
				g.Attr("onchange", "submitTableForm(this.form)"),
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

		// page size select
		h.Select(
			h.Class("select page-size-select"),
			h.Name(p.PageSizeQueryKey),
			g.Group(g.Map(
				p.PageSizeOptions,
				func(o int) g.Node {
					oStr := strconv.Itoa(o)
					return h.Option(
						h.Value(oStr),
						g.Text(oStr),
						g.If(o == p.PageSize, h.Selected()),
					)
				})),
			g.Attr("onchange", "updatePageSizeAndSubmit(this)"),
		),
	)
}

type TableProps struct {
	Classes      c.Classes
	Columns      TableColumns
	Rows         TableRows
	Sort         appsort.Sort
	SortQueryKey string
	Pagination   *TablePaginationProps
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
			p.Pagination != nil,
			TablePagination(p.Pagination),
		),
		h.Div(
			h.Class("table-scroll"),
			h.Table(
				tableHead(&tableHeadProps{
					columns:      p.Columns,
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
				p.Pagination != nil,
				h.Input(
					h.Type("radio"),
					h.Checked(),
					h.Name(p.Pagination.CurrentPageQueryKey),
					h.Value(fmt.Sprintf("%d", p.Pagination.CurrentPage)),
					h.StyleAttr("display: none"),
				),
			),
		),
		g.If(
			p.Pagination != nil,
			TablePagination(p.Pagination),
		),
	)
}
