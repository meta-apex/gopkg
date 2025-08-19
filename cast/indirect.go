package cast

import (
	"reflect"
)

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(i any) (any, bool) {
	if i == nil {
		return nil, false
	}

	if t := reflect.TypeOf(i); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return i, false
	}

	v := reflect.ValueOf(i)

	for v.Kind() == reflect.Ptr || (v.Kind() == reflect.Interface && v.Elem().Kind() == reflect.Ptr) {
		if v.IsNil() {
			return nil, true
		}

		v = v.Elem()
	}

	return v.Interface(), true
}
