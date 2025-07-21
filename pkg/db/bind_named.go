package db

import (
	"fmt"
	"strings"
	"unicode"
)

// BindNamed converts a SQL query with named parameters to positional form
// e.g. `SELECT * FROM users WHERE id = :user_id` becomes `SELECT * FROM users WHERE id = $1`
func BindNamed(query string, args map[string]any) (string, []any, error) {
	var (
		sb         strings.Builder
		inSingle   bool
		inDouble   bool
		inComment  bool
		inMultiCom bool
		paramMap   = make(map[string]int)
		paramOrder []string
		i          int
	)

	for i < len(query) {
		c := query[i]

		if !inSingle && !inDouble && !inComment && !inMultiCom {
			if i+1 < len(query) {
				next := query[i+1]
				if c == '-' && next == '-' {
					inComment = true
					sb.WriteByte(c)
					sb.WriteByte(next)
					i += 2
					continue
				} else if c == '/' && next == '*' {
					inMultiCom = true
					sb.WriteByte(c)
					sb.WriteByte(next)
					i += 2
					continue
				}
			}
		}

		if inComment && c == '\n' {
			inComment = false
		}
		if inMultiCom && i+1 < len(query) && c == '*' && query[i+1] == '/' {
			inMultiCom = false
			sb.WriteByte(c)
			sb.WriteByte(query[i+1])
			i += 2
			continue
		}

		if inComment || inMultiCom {
			sb.WriteByte(c)
			i++
			continue
		}

		if c == '\'' && !inDouble {
			inSingle = !inSingle
			sb.WriteByte(c)
			i++
			continue
		}
		if c == '"' && !inSingle {
			inDouble = !inDouble
			sb.WriteByte(c)
			i++
			continue
		}

		if c == ':' && !inSingle && !inDouble {
			if i > 0 && query[i-1] == ':' {
				sb.WriteByte(c)
				i++
				continue
			}

			start := i + 1
			j := start
			for j < len(query) && (unicode.IsLetter(rune(query[j])) || unicode.IsDigit(rune(query[j])) || query[j] == '_') {
				j++
			}
			name := query[start:j]
			if name == "" {
				return "", nil, fmt.Errorf("empty named parameter at position %d", i)
			}

			pos, exists := paramMap[name]
			if !exists {
				paramOrder = append(paramOrder, name)
				pos = len(paramOrder)
				paramMap[name] = pos
			}

			sb.WriteString(fmt.Sprintf("$%d", pos))
			i = j
			continue
		}

		sb.WriteByte(c)
		i++
	}

	finalArgs := make([]any, len(paramOrder))
	for idx, key := range paramOrder {
		val, ok := args[key]
		if !ok {
			return "", nil, fmt.Errorf("missing parameter: %s", key)
		}
		finalArgs[idx] = val
	}

	return sb.String(), finalArgs, nil
}
