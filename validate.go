package phi

import (
	"fmt"
	"reflect"
	"strings"
)

// validates the data structure
func handleResolve(prefix string, errs *[]string, data interface{}) error {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Struct:
		if err := handleStruct(prefix, errs, data); err != nil {
			return err
		}
	case reflect.Slice, reflect.Array:
		if err := handleSliceArray(prefix, errs, data); err != nil {
			return err
		}
	case reflect.Map:
		if err := handleMap(prefix, errs, data); err != nil {
			return err
		}
	case reflect.Ptr:
		val := reflect.ValueOf(data)
		if val.IsNil() {
			return nil
		}

		return handleResolve(prefix, errs, val.Elem().Interface())
	}

	return nil
}

// handle struct type validation
func handleStruct(prefix string, errs *[]string, data interface{}) error {
	fields := reflect.ValueOf(data)
	for i := 0; i < fields.NumField(); i++ {
		jsonTags := fields.Type().Field(i).Tag.Get("json")

		tag := strings.Split(jsonTags, ",")[0]
		if prefix != "" {
			tag = prefix + "." + tag
		}

		// check if tags contain required and the value is zero
		if strings.Contains(jsonTags, "required") && fields.Field(i).IsZero() {
			*errs = append(*errs, tag)
		} else {
			// no error, validate underlying type(s)
			if err := handleResolve(tag, errs, fields.Field(i).Interface()); err != nil {
				return err
			}
		}
	}

	return nil
}

// handle slice and array type validation
func handleSliceArray(prefix string, errs *[]string, data interface{}) error {
	// get the type of the slice/array and resolve pointer
	v := reflect.ValueOf(data)

	// get kind of the element type
	kind := v.Kind()
	if kind != reflect.Struct && kind != reflect.Slice && kind != reflect.Array && kind != reflect.Ptr && kind != reflect.Map {
		return nil
	}

	// iterate through the slice/array
	for i := 0; i < v.Len(); i++ {
		if err := handleResolve(prefix+fmt.Sprintf("[%d]", i), errs, v.Index(i).Interface()); err != nil {
			return err
		}
	}

	return nil
}

// handle map type validation
func handleMap(prefix string, errs *[]string, data interface{}) error {
	v := reflect.ValueOf(data)

	for _, key := range v.MapKeys() {
		// get the kind of the key and resolve pointer values
		val := v.MapIndex(key)
		if err := handleResolve(prefix+fmt.Sprintf("[%v]", key.Interface()), errs, val.Interface()); err != nil {
			return err
		}
	}

	return nil
}
