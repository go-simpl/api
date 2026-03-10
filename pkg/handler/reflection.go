package handler

// This file binds request data to handler input structs using struct tags: body (json/urlencoded/multipart), path, query, header, cookie.

import (
	"fmt"
	"mime/multipart"
	"reflect"
	"strconv"

	"github.com/go-simpl/simplapi/errors"
	"github.com/go-simpl/simplapi/pkg/framework"
)

// computeValuesFromRequest fills pVal from the request. Tag precedence: body first, then nested struct, then path/query/header/cookie.
func computeValuesFromRequest(f framework.FrameworkRequest, pType reflect.Type, pVal reflect.Value) error {
	for i := 0; i < pVal.NumField(); i++ {
		if pType.Field(i).Tag.Get("body") == "json" {
			b := pVal.Field(i).Addr().Interface()
			err := f.ParseJSONBody(b)
			if err != nil {
				return errors.TypeError{
					FieldName: pType.Field(i).Name,
					Message:   err.Error(),
				}
			}

		} else if pType.Field(i).Tag.Get("body") == "urlencoded" {
			for j := 0; j < pVal.Field(i).NumField(); j++ {
				err := setValue(pVal.Field(i).Field(j), f.GetFormValue(pType.Field(i).Type.Field(j).Tag.Get("form")), pType.Field(i).Type.Field(j).Name)
				if err != nil {
					return err
				}
			}
		} else if pType.Field(i).Tag.Get("body") == "multipart" {
			for j := 0; j < pVal.Field(i).NumField(); j++ {
				fileType := reflect.TypeOf((*multipart.FileHeader)(nil))

				if pVal.Field(i).Field(j).Type().ConvertibleTo(fileType) {
					file, err := f.GetFile(pType.Field(i).Type.Field(j).Tag.Get("form"))
					if err != nil {
						return errors.TypeError{
							FieldName: pType.Field(i).Type.Field(j).Name,
							Message:   err.Error(),
						}
					}
					pVal.Field(i).Field(j).Set(reflect.ValueOf(file))

				} else {
					err := setValue(pVal.Field(i).Field(j), f.GetFormValue(pType.Field(i).Type.Field(j).Tag.Get("form")), pType.Field(i).Type.Field(j).Name)
					if err != nil {
						return err
					}
				}
			}

		} else if pType.Field(i).Type.Kind() == reflect.Struct {
			err := computeValuesFromRequest(f, pType.Field(i).Type, pVal.Field(i))
			if err != nil {
				return err
			}

		} else if pType.Field(i).Tag.Get("path") != "" {
			err := setValue(pVal.Field(i), f.GetPathParam(pType.Field(i).Tag.Get("path")), pType.Field(i).Name)
			if err != nil {
				return err
			}

		} else if pType.Field(i).Tag.Get("query") != "" {
			// Optional (ptr) param: empty string leaves field as zero value.
			if pVal.Field(i).Type().Kind() != reflect.Ptr || f.GetQueryParam(pType.Field(i).Tag.Get("query")) != "" {
				err := setValue(pVal.Field(i), f.GetQueryParam(pType.Field(i).Tag.Get("query")), pType.Field(i).Name)
				if err != nil {
					return err
				}
			}

		} else if pType.Field(i).Tag.Get("header") != "" {
			// Optional (ptr) param: empty string leaves field as zero value.
			if pVal.Field(i).Type().Kind() != reflect.Ptr || f.GetHeader(pType.Field(i).Tag.Get("header")) != "" {
				err := setValue(pVal.Field(i), f.GetHeader(pType.Field(i).Tag.Get("header")), pType.Field(i).Name)
				if err != nil {
					return err
				}
			}
		} else if pType.Field(i).Tag.Get("cookie") != "" {
			// Optional (ptr) param: empty string leaves field as zero value.
			if pVal.Field(i).Type().Kind() != reflect.Ptr || f.GetCookieValue(pType.Field(i).Tag.Get("cookie")) != "" {
				err := setValue(pVal.Field(i), f.GetCookieValue(pType.Field(i).Tag.Get("cookie")), pType.Field(i).Name)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// setValue sets valueObj from the string value. Ptr type means optional: missing or empty string leaves nil/zero.
func setValue(valueObj reflect.Value, value string, fieldName string) error {
	required := true

	if valueObj.Kind() == reflect.Ptr {
		if value == "" {
			return nil
		}

		valueObj.Set(reflect.New(valueObj.Type().Elem()))
		valueObj = valueObj.Elem()
		required = false
	}

	if required && value == "" {
		return errors.TypeError{
			FieldName: fieldName,
			Message:   "value is required",
		}
	}

	switch valueObj.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intValue, err := strconv.ParseInt(value, 10, valueObj.Type().Bits())
		if err != nil {
			return errors.TypeError{
				FieldName: fieldName,
				Message:   err.Error(),
			}
		}
		valueObj.SetInt(intValue)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintValue, err := strconv.ParseUint(value, 10, valueObj.Type().Bits())
		if err != nil {
			return errors.TypeError{
				FieldName: fieldName,
				Message:   err.Error(),
			}
		}
		valueObj.SetUint(uintValue)

	case reflect.Float32, reflect.Float64:
		floatValue, err := strconv.ParseFloat(value, valueObj.Type().Bits())
		if err != nil {
			return errors.TypeError{
				FieldName: fieldName,
				Message:   err.Error(),
			}
		}
		valueObj.SetFloat(floatValue)

	case reflect.Bool:
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return errors.TypeError{
				FieldName: fieldName,
				Message:   err.Error(),
			}
		}
		valueObj.SetBool(boolValue)

	case reflect.String:
		valueObj.SetString(value)

	default:
		return fmt.Errorf("unsupported type: %v", valueObj.Kind())
	}

	return nil
}
