package validation

import (
	"fmt"
	"net/mail"
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

func plural(n int) string {
	if n != 1 {
		return "s"
	}
	return ""
}

func MinLength(s string, min int, ve *ValidationErrors, key string) {
	if len(s) < min {
		ve.Add(key, fmt.Sprintf("must be at least %d character%s", min, plural(min)))
	}
}

func MaxLength(s string, max int, ve *ValidationErrors, key string) {
	if len(s) > max {
		ve.Add(key, fmt.Sprintf("must be at most %d character%s", max, plural(max)))
	}
}

func Lowercase(s string, ve *ValidationErrors, key string) {
	if s != "" && s != strings.ToLower(s) {
		ve.Add(key, "must be all lowercase")
	}
}

func Uppercase(s string, ve *ValidationErrors, key string) {
	if s != "" && s != strings.ToUpper(s) {
		ve.Add(key, "must be all uppercase")
	}
}

func Email(s string, ve *ValidationErrors, key string) {
	_, err := mail.ParseAddress(s)
	if err != nil {
		ve.Add(key, "must be a valid email address")
	}
}

// Password must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, one digit, and one special character
var specialChars = []rune{'!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '-', '_', '+', '=', '{', '}', '[', ']', '|', '\\', ':', ';', '"', '\'', '<', '>', ',', '.', '?', '/'}

func Password(s string, ve *ValidationErrors, key string) {
	errMsg := "must be at least 8 characters long, contain at least one uppercase letter, one lowercase letter, one digit, and one special character"
	hasError := false
	if len(s) < 8 {
		hasError = true
	} else if strings.ToLower(s) == s {
		hasError = true
	} else if strings.ToUpper(s) == s {
		hasError = true
	} else if strings.IndexFunc(s, func(r rune) bool { return '0' <= r && r <= '9' }) == -1 {
		hasError = true
	} else if strings.IndexFunc(s, func(r rune) bool {
		for _, c := range specialChars {
			if r == c {
				return true
			}
		}
		return false
	}) == -1 {
		hasError = true
	}

	if hasError {
		ve.Add(key, errMsg)
	}
}
