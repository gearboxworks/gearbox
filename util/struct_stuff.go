package util

import (
	"reflect"
)

func StructValue(s interface{}) reflect.Value {
	v := reflect.ValueOf(s)

	// if pointer get the underlying element≤
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic("not struct")
	}

	return v
}

func StructFields(s interface{}, tagname string) []reflect.StructField {
	t := StructValue(s).Type()

	var f []reflect.StructField

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// no private fields
		if field.PkgPath != "" {
			continue
		}

		// don't check if it's omitted
		if tag := field.Tag.Get(tagname); tag == "-" {
			continue
		}

		f = append(f, field)
	}
	return f
}

func StructMap(m interface{}, tagname string) map[string]interface{} {
	v := StructValue(m)
	fs := StructFields(m, "json")
	msi := make(map[string]interface{}, len(fs))
	for i, f := range fs {
		tag := f.Tag.Get(tagname)
		if tag == "" {
			continue
		}
		msi[tag] = v.Field(i).Interface()
	}
	return msi
}
