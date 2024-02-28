package utils

import (
	"database/sql"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// settable field kinds
var settableKinds = []reflect.Kind{
	reflect.String,
	reflect.Bool,
	reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
	reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
	reflect.Float32, reflect.Float64,
}

// settable struct types
var settableStructTypes = []reflect.Type{
	reflect.TypeOf(time.Time{}),
	reflect.TypeOf(sql.NullString{}),
	reflect.TypeOf(sql.NullBool{}),
	reflect.TypeOf(sql.NullInt64{}),
	reflect.TypeOf(sql.NullInt32{}),
	reflect.TypeOf(sql.NullInt16{}),
	reflect.TypeOf(sql.NullFloat64{}),
	reflect.TypeOf(sql.NullTime{}),
}

func isSettableField(field reflect.Value) bool {
	fieldKind := field.Kind()
	for _, kind := range settableKinds {
		if fieldKind == kind {
			return true
		}
	}
	if fieldKind == reflect.Struct {
		for _, t := range settableStructTypes {
			if field.Type() == t {
				return true
			}
		}
	}
	return false
}

func parseTime(strValue string) (time.Time, error) {
	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"02/01/2006", // assuming DD/MM/YYYY format
		"02/01/2006 15:04:05",
		// Add more formats as needed
	}

	var parsedValue time.Time
	var err error

	for _, format := range formats {
		parsedValue, err = time.Parse(format, strValue)
		if err == nil {
			// Successfully parsed, return the result
			return parsedValue, nil
		}
	}

	// If none of the formats worked, return an error
	return time.Time{}, fmt.Errorf("could not parse date: %s", strValue)
}

func setField(field reflect.Value, strValue string) error {
	if !field.IsValid() {
		return fmt.Errorf("field is not valid")
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(strValue)
	case reflect.Bool:
		if strValue == "" {
			field.SetBool(false)
		} else {
			parsedValue, err := strconv.ParseBool(strValue)
			if err != nil {
				return err
			}
			field.SetBool(parsedValue)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsedValue, err := strconv.ParseInt(strValue, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(parsedValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsedValue, err := strconv.ParseUint(strValue, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(parsedValue)
	case reflect.Float32, reflect.Float64:
		parsedValue, err := strconv.ParseFloat(strValue, 64)
		if err != nil {
			return err
		}
		field.SetFloat(parsedValue)
	case reflect.Struct:
		if field.Type() == reflect.TypeOf(time.Time{}) {
			parsedValue, err := parseTime(strValue)
			if err != nil {
				return err
			}
			field.Set(reflect.ValueOf(parsedValue))
		} else if field.Type() == reflect.TypeOf(sql.NullString{}) {
			nullString := sql.NullString{String: strValue, Valid: strValue != ""}
			field.Set(reflect.ValueOf(nullString))
		} else if field.Type() == reflect.TypeOf(sql.NullBool{}) {
			parsedValue, err := strconv.ParseBool(strValue)
			nullBool := sql.NullBool{}
			if err == nil {
				nullBool.Bool = parsedValue
				nullBool.Valid = true
			}
			field.Set(reflect.ValueOf(nullBool))
		} else if field.Type() == reflect.TypeOf(sql.NullInt64{}) {
			parsedValue, err := strconv.ParseInt(strValue, 10, 64)
			nullInt64 := sql.NullInt64{}
			if err == nil {
				nullInt64.Int64 = parsedValue
				nullInt64.Valid = true
			}
			field.Set(reflect.ValueOf(nullInt64))
		} else if field.Type() == reflect.TypeOf(sql.NullInt32{}) {
			parsedValue, err := strconv.ParseInt(strValue, 10, 32)
			nullInt32 := sql.NullInt32{}
			if err == nil {
				nullInt32.Int32 = int32(parsedValue)
				nullInt32.Valid = true
			}
			field.Set(reflect.ValueOf(nullInt32))
		} else if field.Type() == reflect.TypeOf(sql.NullInt16{}) {
			parsedValue, err := strconv.ParseInt(strValue, 10, 16)
			nullInt16 := sql.NullInt16{}
			if err == nil {
				nullInt16.Int16 = int16(parsedValue)
				nullInt16.Valid = true
			}
			field.Set(reflect.ValueOf(nullInt16))
		} else if field.Type() == reflect.TypeOf(sql.NullFloat64{}) {
			parsedValue, err := strconv.ParseFloat(strValue, 64)
			nullFloat64 := sql.NullFloat64{}
			if err == nil {
				nullFloat64.Float64 = parsedValue
				nullFloat64.Valid = true
			}
			field.Set(reflect.ValueOf(nullFloat64))
		} else if field.Type() == reflect.TypeOf(sql.NullTime{}) {
			parsedValue, err := parseTime(strValue)
			nullTime := sql.NullTime{}
			if err == nil {
				nullTime.Time = parsedValue
				nullTime.Valid = true
			}
			field.Set(reflect.ValueOf(nullTime))
		}

	}

	return nil
}

func UnmarshalUrlValues(urlValues url.Values, v interface{}) error {
	val := reflect.ValueOf(v).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		fTypeKind := field.Type.Kind()
		fType := field.Type
		fName := field.Name

		// check if the field is settable
		isSettable := isSettableField(fieldValue)

		if isSettable {
			strValue := urlValues.Get(fName)
			err := setField(fieldValue, strValue)
			if err != nil {
				return err
			}
		} else if field.Type.Kind() == reflect.Slice {
			// slice field
			// we support square bracket/index notation and multiple form values under
			// same name

			// filter the form data based on the field name prefix
			filteredForm := make(url.Values)
			for k, v := range urlValues {
				if strings.HasPrefix(k, fName) {
					filteredForm[fName] = v
				}
			}

			// is square bracket notation
			isSquareBracketNotation := false
			for k := range filteredForm {
				if strings.HasPrefix(k, fName+"[") {
					isSquareBracketNotation = true
					break
				}
			}

			if !isSquareBracketNotation {
				// not using square bracket notation, the slice values are all the values
				// under the field name key
				strValues := filteredForm[fName]
				slice := reflect.MakeSlice(field.Type, len(strValues), len(strValues))
				fieldValue.Set(slice)
				elemType := fType.Elem()
				// Populate the slice
				for i, strValue := range strValues {
					subValue := reflect.New(elemType).Elem()
					err := setField(subValue, strValue)
					if err != nil {
						return err
					}
					fieldValue.Index(i).Set(subValue)
				}
			} else {
				// get the indices of the slice by removing the common prefix
				// and then removing everything from the first ']' (inclusive)
				indices := make([]int, 0)
				for k := range filteredForm {
					splitIndex := strings.Index(k, "]")
					if splitIndex == -1 {
						return fmt.Errorf("invalid key: %s", k)
					}
					indexStr := k[len(fName)+1 : splitIndex]
					index, err := strconv.Atoi(indexStr)
					if err != nil {
						return fmt.Errorf("invalid index parsing array form field: %s", indexStr)
					}
					indices = append(indices, index)
				}

				// Ensure the slice has enough capacity, use the max index
				sliceLen := 0
				for _, index := range indices {
					if index > sliceLen-1 {
						sliceLen = index + 1
					}
				}

				slice := reflect.MakeSlice(field.Type, sliceLen, sliceLen)
				fieldValue.Set(slice)
				// Populate the slice - we can use the indices data and call DecodeForm
				// for each index
				elemType := fType.Elem()
				for _, index := range indices {
					subForm := make(url.Values)
					prefix := fmt.Sprintf("%s[%d]", fName, index)
					for k, v := range urlValues {
						if strings.HasPrefix(k, prefix) {
							subForm[strings.TrimPrefix(k, prefix)] = v
						}
					}
					subValue := reflect.New(elemType).Elem()
					if err := UnmarshalUrlValues(subForm, subValue.Addr().Interface()); err != nil {
						return err
					}
					fieldValue.Index(index).Set(subValue)
				}
			}

		} else if field.Type.Kind() == reflect.Struct {
			// struct field - we recurse
			subForm := make(url.Values)
			prefix := fmt.Sprintf("%s.", fName)
			for k, v := range urlValues {
				if strings.HasPrefix(k, prefix) {
					subForm[strings.TrimPrefix(k, prefix)] = v
				}
			}
			if err := UnmarshalUrlValues(subForm, fieldValue.Addr().Interface()); err != nil {
				return err
			}
		} else {
			// unknown field
			return fmt.Errorf("unknown field type kind: %s and field type: %s", fTypeKind, fType)
		}
	}

	return nil
}
