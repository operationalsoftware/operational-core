package modelutil

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// Known acronyms we want to handle as units
var knownAcronyms = map[string]string{
	"URL":  "Url",
	"ID":   "Id",
	"API":  "Api",
	"HTTP": "Http",
}

// Preprocess acronyms in the string
func preprocessAcronyms(str string) string {
	for k, v := range knownAcronyms {
		str = strings.ReplaceAll(str, k, v)
	}
	return str
}

// Regex patterns for CamelCase to snake_case
var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// Convert CamelCase to snake_case, accounting for acronyms
func toSnakeCase(str string) string {
	str = preprocessAcronyms(str)
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// GetFieldColumnName returns the `column` tag if present, else snake_case version of fieldName
func GetFieldColumnName(modelType any, fieldName string) (string, error) {
	t := reflect.TypeOf(modelType)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("provided value is not a struct or pointer to struct")
	}

	field, found := t.FieldByName(fieldName)
	if !found {
		return "", fmt.Errorf("field %s not found", fieldName)
	}

	tagValue := field.Tag.Get("column")
	if tagValue != "" {
		return tagValue, nil
	}

	return toSnakeCase(fieldName), nil
}

func IsFieldSortable(modelType any, fieldName string) (bool, error) {
	t := reflect.TypeOf(modelType)

	// If pointer, dereference
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Only operate on structs
	if t.Kind() != reflect.Struct {
		return false, fmt.Errorf("provided value is not a struct or pointer to struct")
	}

	// If it's a pointer, get the underlying element
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	field, found := t.FieldByName(fieldName)
	if !found {
		return false, fmt.Errorf("field %s not found", fieldName)
	}

	tagValue := field.Tag.Get("sortable")
	return tagValue == "true", nil
}
