package defaults

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

var (
	timeDurationType = reflect.TypeOf(time.Second)
)

// error list
var (
	// ErrPassValue is returned if the caller pass a value instead of a pointer
	ErrPassValue = errors.New("must pass a pointer, not a value")
	// ErrNotStruct is returned if the caller pass a pointer of non struct
	ErrNotStruct = errors.New("must pass a pointer of struct")
)

// SetDefault set default value from struct tag: default
// for example:
// type A struct {
//		S string `default:"this is default"`
// }
func SetDefault(v interface{}) error {
	val := reflect.ValueOf(v)
	// prevent silent error, if a value is sent, the original value won't change
	if val.Kind() != reflect.Ptr {
		return ErrPassValue
	}
	indirect := reflect.Indirect(val)
	// prevent panic when call NumField()
	if indirect.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	numfield := indirect.NumField()
	for i := 0; i < numfield; i++ {
		fi := indirect.Field(i)
		if !fi.CanSet() {
			continue
		}

		// continue if it is not empty value
		if !reflect.DeepEqual(reflect.Zero(fi.Type()).Interface(), fi.Interface()) {
			continue
		}

		f := indirect.Type().Field(i)
		t := f.Tag.Get("default")
		// continue if default tag is not available
		if t == "" {
			continue
		}

		// for special types which have their own parser
		switch f.Type {
		case timeDurationType:
			n, err := time.ParseDuration(t)
			if err != nil {
				return err
			}
			fi.Set(reflect.ValueOf(n))
			continue
		}

		// for primitive types
		switch f.Type.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			n, err := strconv.ParseInt(t, 10, 64)
			if err != nil {
				return err
			}
			fi.Set(reflect.ValueOf(n).Convert(f.Type))
		case reflect.Float32, reflect.Float64:
			n, err := strconv.ParseFloat(t, 64)
			if err != nil {
				return err
			}
			fi.Set(reflect.ValueOf(n).Convert(f.Type))
		case reflect.String:
			fi.Set(reflect.ValueOf(t).Convert(f.Type))
		}
	}
	return nil
}

// ReplaceDefaultFrom will replace value of a struct with another value in a struct
// to be replaced, field name and type between two struct must be the same
// for example:
// type A struct {
//		Field1 string
// }
//
// type B struct {
// 		Field1 string
//}
func ReplaceDefaultFrom(source interface{}, replacer interface{}) error {
	val := reflect.ValueOf(source)
	// prevent silent error, if a value is sent, the original value won't change
	if val.Kind() != reflect.Ptr {
		return ErrPassValue
	}

	inreplacer := reflect.Indirect(reflect.ValueOf(replacer))
	insource := reflect.Indirect(val)
	// prevent panic when call NumField() and FieldByName()
	if insource.Kind() != reflect.Struct || inreplacer.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	numfield := insource.NumField()
	for i := 0; i < numfield; i++ {
		field := insource.Field(i)
		fieldName := insource.Type().Field(i).Name
		replacerfield := inreplacer.FieldByName(fieldName)

		// skip if cannot set field or field not found in replacer or have different kind
		if !field.CanSet() || !replacerfield.IsValid() || field.Kind() != replacerfield.Kind() {
			continue
		}

		// continue if it is not empty value
		if !reflect.DeepEqual(reflect.Zero(field.Type()).Interface(), field.Interface()) {
			continue
		}

		field.Set(replacerfield)
	}
	return nil
}
