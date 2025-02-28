package nilsafe

import "time"

// Str returns the string value if not nil, otherwise returns an empty string.
func Str(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Int returns the int value if not nil, otherwise returns 0.
func Int(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// Int64 returns the int64 value if not nil, otherwise returns 0.
func Int64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

// Float32 returns the float32 value if not nil, otherwise returns 0.0.
func Float32(f *float32) float32 {
	if f == nil {
		return 0.0
	}
	return *f
}

// Float64 returns the float64 value if not nil, otherwise returns 0.0.
func Float64(f *float64) float64 {
	if f == nil {
		return 0.0
	}
	return *f
}

// Bool returns the bool value if not nil, otherwise returns false.
func Bool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// Time returns the time value if not nil, otherwise returns zero time.
func Time(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

// Bytes returns the byte slice if not nil, otherwise returns an empty slice.
func Bytes(b *[]byte) []byte {
	if b == nil {
		return []byte{}
	}
	return *b
}
