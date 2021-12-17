package query

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
)

func Marshal(s interface{}, values url.Values) error {
	val := reflect.ValueOf(s)
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return ErrStruct
		}

		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return ErrStruct
	}

	return reflectValueFromTag4Marshal(val, values)
}

func reflectSliceValueFromTag4Marshal(tag string, val reflect.Value, values url.Values) {
	switch t := val.Type().Elem().Kind(); t {
	case reflect.String:
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i)
			values.Add(tag, v.String())
		}
	case reflect.Bool:
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i)
			values.Add(tag, strconv.FormatBool(v.Bool()))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i)
			values.Add(tag, strconv.FormatUint(v.Uint(), 10))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i)
			values.Add(tag, strconv.FormatInt(v.Int(), 10))
		}
	case reflect.Float32, reflect.Float64:
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i)
			values.Add(tag, strconv.FormatFloat(v.Float(), 'f', 8, 64))
		}
	}
}

// nolint: gocognit
func reflectValueFromTag4Marshal(val reflect.Value, values url.Values) error {
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		kt := typ.Field(i)
		tagQs := strings.Split(kt.Tag.Get("query"), ",")
		tag := strings.Trim(tagQs[0], " ")

		var omitEmpty bool

		for idx := 1; idx < len(tagQs); idx++ {
			if strings.Trim(tagQs[idx], " ") == "omitempty" {
				omitEmpty = true

				break
			}
		}

		if tag == "-" {
			continue
		}

		tag, _ = parseTag(tag)

		if tag == "" {
			tag = strcase.ToSnake(kt.Name)
		}

		sv := val.Field(i)

		switch sv.Kind() {
		case reflect.String:
			if omitEmpty && sv.String() == "" {
				continue
			}

			values.Add(tag, sv.String())
		case reflect.Bool:
			if omitEmpty && !sv.Bool() {
				continue
			}

			values.Add(tag, strconv.FormatBool(sv.Bool()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			if omitEmpty && sv.Uint() == 0 {
				continue
			}

			values.Add(tag, strconv.FormatUint(sv.Uint(), 10))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if omitEmpty && sv.Int() == 0 {
				continue
			}

			values.Add(tag, strconv.FormatInt(sv.Int(), 10))
		case reflect.Float32, reflect.Float64:
			if omitEmpty && sv.Float() == 0 {
				continue
			}

			values.Add(tag, strconv.FormatFloat(sv.Float(), 'f', 8, 64))
		case reflect.Struct:
			err := reflectValueFromTag4Marshal(sv, values)
			if err != nil {
				return newParseError(err, "cast from struct", tag, sv.Type().String(), "", tag)
			}
		case reflect.Ptr:
			if sv.IsNil() {
				if omitEmpty {
					continue
				}

				sv.Set(reflect.New(sv.Type().Elem()))
			}

			err := reflectValueFromTag4Marshal(sv.Elem(), values)

			if err != nil {
				return newParseError(err, "cast from ptr", tag, sv.Type().String(), "", tag)
			}
		case reflect.Slice:
			if sv.Len() <= 0 {
				continue
			}

			reflectSliceValueFromTag4Marshal(tag, sv, values)
		}
	}

	return nil
}
