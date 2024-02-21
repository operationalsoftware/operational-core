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

func DecodeForm(form url.Values, v interface{}) error {
	val := reflect.ValueOf(v).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i)

		fTypeKind := field.Type.Kind()
		fType := field.Type
		fName := field.Name

		if fTypeKind == reflect.String {
			// string field
			fieldValue.SetString(form.Get(fName))
		} else if fTypeKind == reflect.Bool {
			// bool field
			var value bool
			if form.Get(fName) == "" {
				// the field is not present in the from e.g. unchecked checkbox
				value = false
			} else {
				parsedValue, err := strconv.ParseBool(form.Get(fName))
				if err != nil {
					return err
				}
				value = parsedValue
			}
			fieldValue.SetBool(value)
		} else if fTypeKind == reflect.Int ||
			fTypeKind == reflect.Int8 ||
			fTypeKind == reflect.Int16 ||
			fTypeKind == reflect.Int32 ||
			fTypeKind == reflect.Int64 {
			// int field
			value, err := strconv.ParseInt(form.Get(fName), 10, 64)
			if err != nil {
				return err
			}
			fieldValue.SetInt(value)
		} else if fTypeKind == reflect.Uint ||
			fTypeKind == reflect.Uint8 ||
			fTypeKind == reflect.Uint16 ||
			fTypeKind == reflect.Uint32 ||
			fTypeKind == reflect.Uint64 {
			// uint field
			value, err := strconv.ParseUint(form.Get(fName), 10, 64)
			if err != nil {
				return err
			}
			fieldValue.SetUint(value)
		} else if fTypeKind == reflect.Float32 ||
			fTypeKind == reflect.Float64 {
			// float field
			value, err := strconv.ParseFloat(form.Get(fName), 64)
			if err != nil {
				return err
			}
			fieldValue.SetFloat(value)
		} else if fTypeKind == reflect.Struct && fType == reflect.TypeOf(time.Time{}) {
			// time field
			timeValue, err := time.Parse(time.RFC3339, form.Get(fName))
			if err != nil {
				return err
			}
			fieldValue.Set(reflect.ValueOf(timeValue))
		} else if fTypeKind == reflect.Struct && fType == reflect.TypeOf(sql.NullString{}) {
			// sql.NullString field
			strValue := form.Get(fName)
			nullString := sql.NullString{String: strValue, Valid: strValue != ""}
			fieldValue.Set(reflect.ValueOf(nullString))
		} else if fTypeKind == reflect.Struct && fType == reflect.TypeOf(sql.NullBool{}) {
			// sql.NullBool field
			boolValue, err := strconv.ParseBool(form.Get(fName))
			nullBool := sql.NullBool{}
			if err == nil {
				nullBool.Bool = boolValue
				nullBool.Valid = true
			}
			fieldValue.Set(reflect.ValueOf(nullBool))
		} else if fTypeKind == reflect.Struct && fType == reflect.TypeOf(sql.NullInt64{}) {
			// sql.NullInt64 field
			intValue, err := strconv.ParseInt(form.Get(fName), 10, 64)
			nullInt64 := sql.NullInt64{}
			if err == nil {
				nullInt64.Int64 = intValue
				nullInt64.Valid = true
			}
			fieldValue.Set(reflect.ValueOf(nullInt64))
		} else if fTypeKind == reflect.Struct && fType == reflect.TypeOf(sql.NullInt32{}) {
			// sql.NullInt32 field
			intValue, err := strconv.ParseInt(form.Get(fName), 10, 32)
			nullInt32 := sql.NullInt32{}
			if err == nil {
				nullInt32.Int32 = int32(intValue)
				nullInt32.Valid = true
			}
			fieldValue.Set(reflect.ValueOf(nullInt32))
		} else if fTypeKind == reflect.Struct && fType == reflect.TypeOf(sql.NullInt16{}) {
			// sql.NullInt16 field
			intValue, err := strconv.ParseInt(form.Get(fName), 10, 16)
			nullInt16 := sql.NullInt16{}
			if err == nil {
				nullInt16.Int16 = int16(intValue)
				nullInt16.Valid = true
			}
			fieldValue.Set(reflect.ValueOf(nullInt16))
		} else if fTypeKind == reflect.Struct && fType == reflect.TypeOf(sql.NullFloat64{}) {
			// sql.NullFloat64 field
			floatValue, err := strconv.ParseFloat(form.Get(fName), 64)
			nullFloat64 := sql.NullFloat64{}
			if err == nil {
				nullFloat64.Float64 = floatValue
				nullFloat64.Valid = true
			}
			fieldValue.Set(reflect.ValueOf(nullFloat64))
		} else if fTypeKind == reflect.Struct && fType == reflect.TypeOf(sql.NullTime{}) {
			// sql.NullTime field
			timeValue, err := time.Parse(time.RFC3339, form.Get(fName))
			nullTime := sql.NullTime{}
			if err == nil {
				nullTime.Time = timeValue
				nullTime.Valid = true
			}
			fieldValue.Set(reflect.ValueOf(nullTime))
		} else if fTypeKind == reflect.Slice {
			// slice field
			elemType := field.Type.Elem()
			commonPrefix := fmt.Sprintf("%s[", fName)
			// filter the form data based on this prefix
			filteredForm := make(url.Values)
			for k, v := range form {
				if strings.HasPrefix(k, commonPrefix) {
					filteredForm[k] = v
				}
			}
			// get the indices of the slice by removing the common prefix
			// and then removing everything after the first ']' (inclusive)
			indices := make([]int, 0)
			for k := range filteredForm {
				splitIndex := strings.Index(k, "]")
				if splitIndex == -1 {
					return fmt.Errorf("invalid key: %s", k)
				}
				indexStr := k[len(commonPrefix):splitIndex]
				index, err := strconv.Atoi(indexStr)
				if err != nil {
					return fmt.Errorf("invalid index parsing array form field: %s", indexStr)
				}
				indices = append(indices, index)
			}
			// Ensure the slice has enough capacity, use the max index
			maxIndex := 0
			for _, index := range indices {
				if index > maxIndex {
					maxIndex = index
				}
			}
			slice := reflect.MakeSlice(field.Type, maxIndex+1, maxIndex+1)
			fieldValue.Set(slice)
			// Populate the slice - we can use the indices data and call DecodeForm
			// for each index
			for _, index := range indices {
				subForm := make(url.Values)
				prefix := fmt.Sprintf("%s[%d]", fName, index)
				for k, v := range form {
					if strings.HasPrefix(k, prefix) {
						subForm[strings.TrimPrefix(k, prefix)] = v
					}
				}
				subValue := reflect.New(elemType).Elem()
				if err := DecodeForm(subForm, subValue.Addr().Interface()); err != nil {
					return err
				}
				fieldValue.Index(index).Set(subValue)
			}
		} else if fTypeKind == reflect.Struct {
			// struct field
			subForm := make(url.Values)
			prefix := fmt.Sprintf("%s.", fName)
			for k, v := range form {
				if strings.HasPrefix(k, prefix) {
					subForm[strings.TrimPrefix(k, prefix)] = v
				}
			}
			if err := DecodeForm(subForm, fieldValue.Addr().Interface()); err != nil {
				return err
			}
		} else {
			// unknown field
			return fmt.Errorf("unknown field type kind: %s and field type: %s", fTypeKind, fType)
		}
	}

	return nil
}
