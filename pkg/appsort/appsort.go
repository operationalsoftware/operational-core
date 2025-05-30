package appsort

import (
	"app/pkg/modelutil"
	"fmt"
	"strings"
)

type Direction string

const (
	DirectionAsc  Direction = "asc"
	DirectionDesc Direction = "desc"
)

type SortItem struct {
	Field string
	Sort  Direction
}

type Sort []SortItem

func (s *Sort) ParseQueryParam(modelType any, queryParam string) error {
	sort := Sort{}
	if queryParam != "" {
		sortItems := strings.Split(queryParam, "-")
		for i := 0; i < len(sortItems); i += 2 {
			field := sortItems[i]
			direction := Direction(sortItems[i+1])
			// validate direction
			if direction != DirectionAsc && direction != DirectionDesc {
				continue
			}
			// validate sortable field
			isFieldSortable, err := modelutil.IsFieldSortable(modelType, field)
			if err != nil {
				return fmt.Errorf("error parsing sort: %v", err)
			}
			if !isFieldSortable {
				return fmt.Errorf("error parsing sort: %v", err)
			}

			sort = append(sort, SortItem{
				Field: sortItems[i],
				Sort:  direction,
			})
		}
	}

	*s = sort

	return nil
}

func (s *Sort) EncodeQueryParam() string {
	sortString := ""
	for i, si := range *s {
		if i != 0 {
			sortString += "-"
		}
		sortString += si.Field + "-" + string(si.Sort)
	}

	return sortString
}

func (s *Sort) IsSortedBy(key string) bool {
	for _, si := range *s {
		if si.Field == key {
			return true
		}
	}

	return false
}

func (s *Sort) GetIndex(key string) int {
	for i, si := range *s {
		if si.Field == key {
			return i
		}
	}

	return -1
}

func (s *Sort) GetDirection(key string) Direction {
	for _, si := range *s {
		if si.Field == key {
			return si.Sort
		}
	}

	return ""
}

func (s *Sort) ToOrderByClause(modelType any) (string, error) {
	sql := ""
	for i, si := range *s {
		if i != 0 {
			sql += ", "
		}

		isFieldSortable, err := modelutil.IsFieldSortable(modelType, si.Field)
		if !isFieldSortable {
			return "", fmt.Errorf("sort error: field %s is not sortable", si.Field)
		}

		column, err := modelutil.GetFieldColumnName(modelType, si.Field)
		if err != nil {
			return "", fmt.Errorf("sort error: %v", err)
		}

		sql += column + " " + strings.ToUpper(string(si.Sort))
	}

	if sql != "" {
		sql = "ORDER BY " + sql
	}

	return sql, nil
}
