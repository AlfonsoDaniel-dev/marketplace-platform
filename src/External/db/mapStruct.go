package db

import (
	"errors"
	"reflect"
)

func MapStructValues(values []any, object any) error {
	if len(values) == 0 {
		return errors.New("no values to map")
	}

	objectStruct := reflect.ValueOf(object)
	if objectStruct.Kind() != reflect.Ptr || objectStruct.Elem().Kind() != reflect.Struct {
		return errors.New("object is not a pointer to struct")
	}

	for i := 0; i < objectStruct.NumField(); i++ {

		for j := 0; j < len(values); j++ {

			field := objectStruct.Elem().Field(i)

			if !field.CanSet() {
				return errors.New("cannot set field")
			}
			field.Set(reflect.ValueOf(values[j]))
		}
	}

	return nil
}
