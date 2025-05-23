package appsort

import (
	"strings"
)

type Direction string

const (
	DirectionAsc  Direction = "asc"
	DirectionDesc Direction = "desc"
)

type SortItem struct {
	Key  string
	Sort Direction
}

type Sort []SortItem

func (s *Sort) ParseQueryParam(queryParam string, allowedKeys []string) {
	sort := Sort{}
	if queryParam != "" {
		sortItems := strings.Split(queryParam, "-")
		for i := 0; i < len(sortItems); i += 2 {
			key := sortItems[i]
			direction := Direction(sortItems[i+1])
			// validate direction
			if direction != DirectionAsc && direction != DirectionDesc {
				continue
			}
			// validate parsed key is in allowed keys
			for j, allowedKey := range allowedKeys {
				if key == allowedKey {
					break
				} else if j == len(allowedKeys)-1 {
					continue
				}
			}

			sort = append(sort, SortItem{
				Key:  sortItems[i],
				Sort: Direction(sortItems[i+1]),
			})
		}
	}

	*s = sort
}

func (s *Sort) EncodeQueryParam() string {
	sortString := ""
	for i, si := range *s {
		if i != 0 {
			sortString += "-"
		}
		sortString += si.Key + "-" + string(si.Sort)
	}

	return sortString
}

func (s *Sort) IsSortedBy(key string) bool {
	for _, si := range *s {
		if si.Key == key {
			return true
		}
	}

	return false
}

func (s *Sort) GetIndex(key string) int {
	for i, si := range *s {
		if si.Key == key {
			return i
		}
	}

	return -1
}

func (s *Sort) GetDirection(key string) Direction {
	for _, si := range *s {
		if si.Key == key {
			return si.Sort
		}
	}

	return ""
}

func (s *Sort) ToOrderByClause(keyMap map[string]string) string {
	sql := ""
	for i, si := range *s {
		if i != 0 {
			sql += ", "
		}

		column, ok := keyMap[si.Key]
		if !ok {
			column = si.Key
		}

		sql += column + " " + strings.ToUpper(string(si.Sort))
	}

	if sql != "" {
		sql = "ORDER BY " + sql
	}

	return sql
}
