package validate

import (
	"fmt"
	"strings"
)

type ValidationErrors map[string][]string

func (v ValidationErrors) Add(key string, message string) {
	v[key] = append(v[key], message)
}

func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}

// Returns an error message for the given key, or an empty string if the key does not exist
func (v ValidationErrors) GetError(key string, name string) string {
	if v.HasErrors() {
		if errors, ok := v[key]; ok {
			return fmt.Sprintf("%s %s", name, strings.Join(errors, ", "))
		}
	}
	return ""
}
