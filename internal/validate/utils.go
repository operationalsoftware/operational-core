package validate

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/shopspring/decimal"
)

func plural(n int) string {
	if n != 1 {
		return "s"
	}
	return ""
}

func MinLength(ve *ValidationErrors, key string, s string, min int) {
	if len(s) < min {
		ve.Add(key, fmt.Sprintf("must be at least %d character%s", min, plural(min)))
	}
}

func MaxLength(ve *ValidationErrors, key string, s string, max int) {
	if len(s) > max {
		ve.Add(key, fmt.Sprintf("must be at most %d character%s", max, plural(max)))
	}
}

func Lowercase(ve *ValidationErrors, key string, s string) {
	if s != "" && s != strings.ToLower(s) {
		ve.Add(key, "must be all lowercase")
	}
}

func Uppercase(ve *ValidationErrors, key string, s string) {
	if s != "" && s != strings.ToUpper(s) {
		ve.Add(key, "must be all uppercase")
	}
}

func Email(ve *ValidationErrors, key string, s string) {
	_, err := mail.ParseAddress(s)
	if err != nil {
		ve.Add(key, "must be a valid email address")
	}
}

var specialChars = []rune{'!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '-', '_', '+', '=', '{', '}', '[', ']', '|', '\\', ':', ';', '"', '\'', '<', '>', ',', '.', '?', '/'}

func Password(ve *ValidationErrors, key string, s string) {
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

// Decimal comparisons
func DecimalGTE(ve *ValidationErrors, key string, d, min decimal.Decimal) {
	if d.LessThan(min) {
		ve.Add(key, fmt.Sprintf("must be greater than or equal to %s", min.String()))
	}
}

func DecimalGT(ve *ValidationErrors, key string, d, min decimal.Decimal) {
	if d.LessThanOrEqual(min) {
		ve.Add(key, fmt.Sprintf("must be greater than %s", min.String()))
	}
}

func DecimalLTE(ve *ValidationErrors, key string, d, max decimal.Decimal) {
	if d.GreaterThan(max) {
		ve.Add(key, fmt.Sprintf("must be less than or equal to %s", max.String()))
	}
}

func DecimalLT(ve *ValidationErrors, key string, d, max decimal.Decimal) {
	if d.GreaterThanOrEqual(max) {
		ve.Add(key, fmt.Sprintf("must be less than %s", max.String()))
	}
}

// Integer comparisons
func IntGTE(ve *ValidationErrors, key string, value, min int) {
	if value < min {
		ve.Add(key, fmt.Sprintf("must be greater than or equal to %d", min))
	}
}

func IntGT(ve *ValidationErrors, key string, value, min int) {
	if value <= min {
		ve.Add(key, fmt.Sprintf("must be greater than %d", min))
	}
}

func IntLTE(ve *ValidationErrors, key string, value, max int) {
	if value > max {
		ve.Add(key, fmt.Sprintf("must be less than or equal to %d", max))
	}
}

func IntLT(ve *ValidationErrors, key string, value, max int) {
	if value >= max {
		ve.Add(key, fmt.Sprintf("must be less than %d", max))
	}
}

// Unsigned Integer comparisons
func UintGTE(ve *ValidationErrors, key string, value, min uint) {
	if value < min {
		ve.Add(key, fmt.Sprintf("must be greater than or equal to %d", min))
	}
}

func UintGT(ve *ValidationErrors, key string, value, min uint) {
	if value <= min {
		ve.Add(key, fmt.Sprintf("must be greater than %d", min))
	}
}

func UintLTE(ve *ValidationErrors, key string, value, max uint) {
	if value > max {
		ve.Add(key, fmt.Sprintf("must be less than or equal to %d", max))
	}
}

func UintLT(ve *ValidationErrors, key string, value, max uint) {
	if value >= max {
		ve.Add(key, fmt.Sprintf("must be less than %d", max))
	}
}
