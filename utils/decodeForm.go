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

		switch field.Type.Kind() {
		case reflect.String:
			fieldValue.SetString(form.Get(field.Name))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value, err := strconv.ParseInt(form.Get(field.Name), 10, 64)
			if err != nil {
				return err
			}
			fieldValue.SetInt(value)
		case reflect.Slice:
			elemType := field.Type.Elem()
			if elemType.Kind() == reflect.String || elemType.Kind() == reflect.Int || elemType.Kind() == reflect.Int8 ||
				elemType.Kind() == reflect.Int16 || elemType.Kind() == reflect.Int32 || elemType.Kind() == reflect.Int64 {
				keyPrefix := fmt.Sprintf("%s[", field.Name)
				var values []string
				for k, v := range form {
					if strings.HasPrefix(k, keyPrefix) {
						values = append(values, v...)
					}
				}
				slice := reflect.MakeSlice(field.Type, len(values), len(values))
				for j, strValue := range values {
					if elemType.Kind() == reflect.String {
						slice.Index(j).SetString(strValue)
					} else if elemType.Kind() == reflect.Int || elemType.Kind() == reflect.Int8 ||
						elemType.Kind() == reflect.Int16 || elemType.Kind() == reflect.Int32 || elemType.Kind() == reflect.Int64 {
						intValue, err := strconv.ParseInt(strValue, 10, 64)
						if err != nil {
							return err
						}
						slice.Index(j).SetInt(intValue)
					}
				}
				fieldValue.Set(slice)
			} else if elemType.Kind() == reflect.Struct && elemType == reflect.TypeOf(time.Time{}) {
				keyPrefix := fmt.Sprintf("%s[", field.Name)
				var values []string
				for k, v := range form {
					if strings.HasPrefix(k, keyPrefix) {
						values = append(values, v...)
					}
				}
				timeSlice := make([]time.Time, len(values))
				for j, strValue := range values {
					timeValue, err := time.Parse(time.RFC3339, strValue)
					if err != nil {
						return err
					}
					timeSlice[j] = timeValue
				}
				fieldValue.Set(reflect.ValueOf(timeSlice))
			} else if elemType.Kind() == reflect.Struct {
				keyPrefix := fmt.Sprintf("%s[", field.Name)
				var indices []int
				for k := range form {
					if strings.HasPrefix(k, keyPrefix) {
						indexStr := strings.TrimSuffix(strings.TrimPrefix(k, keyPrefix), "]")
						index, err := strconv.Atoi(indexStr)
						if err != nil {
							return err
						}
						indices = append(indices, index)
					}
				}
				// Ensure the slice has enough capacity
				slice := reflect.MakeSlice(field.Type, len(indices), len(indices))
				fieldValue.Set(reflect.AppendSlice(fieldValue, slice))
				// Populate each struct in the slice
				for _, index := range indices {
					subForm := make(url.Values)
					prefix := fmt.Sprintf("%s[%d].", field.Name, index)
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
			}
		case reflect.Struct:
			if field.Type == reflect.TypeOf(time.Time{}) {
				timeValue, err := time.Parse(time.RFC3339, form.Get(field.Name))
				if err != nil {
					return err
				}
				fieldValue.Set(reflect.ValueOf(timeValue))
			} else if field.Type == reflect.TypeOf(sql.NullString{}) {
				strValue := form.Get(field.Name)
				nullString := sql.NullString{String: strValue, Valid: strValue != ""}
				fieldValue.Set(reflect.ValueOf(nullString))
			} else if field.Type == reflect.TypeOf(sql.NullInt64{}) {
				intValue, err := strconv.ParseInt(form.Get(field.Name), 10, 64)
				if err != nil {
					return err
				}
				nullInt64 := sql.NullInt64{Int64: intValue, Valid: true}
				fieldValue.Set(reflect.ValueOf(nullInt64))
			} else if field.Type == reflect.TypeOf(sql.NullInt32{}) {
				intValue, err := strconv.ParseInt(form.Get(field.Name), 10, 32)
				if err != nil {
					return err
				}
				nullInt32 := sql.NullInt32{Int32: int32(intValue), Valid: true}
				fieldValue.Set(reflect.ValueOf(nullInt32))
			} else if field.Type == reflect.TypeOf(sql.NullFloat64{}) {
				floatValue, err := strconv.ParseFloat(form.Get(field.Name), 64)
				if err != nil {
					return err
				}
				nullFloat64 := sql.NullFloat64{Float64: floatValue, Valid: true}
				fieldValue.Set(reflect.ValueOf(nullFloat64))
			} else if field.Type == reflect.TypeOf(sql.NullTime{}) {
				timeValue, err := time.Parse(time.RFC3339, form.Get(field.Name))
				if err != nil {
					return err
				}
				nullTime := sql.NullTime{Time: timeValue, Valid: true}
				fieldValue.Set(reflect.ValueOf(nullTime))
			} else if field.Type == reflect.TypeOf(sql.NullBool{}) {
				boolValue, err := strconv.ParseBool(form.Get(field.Name))
				if err != nil {
					return err
				}
				nullBool := sql.NullBool{Bool: boolValue, Valid: true}
				fieldValue.Set(reflect.ValueOf(nullBool))
			} else {
				subForm := make(url.Values)
				prefix := fmt.Sprintf("%s.", field.Name)
				for k, v := range form {
					if strings.HasPrefix(k, prefix) {
						subForm[strings.TrimPrefix(k, prefix)] = v
					}
				}
				subValue := reflect.New(field.Type).Elem()
				if err := DecodeForm(subForm, subValue.Addr().Interface()); err != nil {
					return err
				}
				fieldValue.Set(subValue)
			}
		}
	}

	return nil
}
