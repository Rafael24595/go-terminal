package table

import "reflect"

type Field struct {
	Header string
	Value  any
}

func StructHeaders[T any]() []string {
	var zero T

	headers := make([]string, 0)
	for _, v := range StructFieds(zero) {
		headers = append(headers, v.Header)
	}

	return headers
}

func StructFieds(s any) []Field {
	if s == nil {
		return make([]Field, 0)
	}

	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	if t.Kind() == reflect.Pointer {
		if v.IsNil() {
			return make([]Field, 0)
		}
		v = v.Elem()
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return make([]Field, 0)
	}

	var result []Field

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()

		result = append(result, Field{
			Header: field.Name,
			Value:  value,
		})
	}

	return result
}
