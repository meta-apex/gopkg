package copier

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type deepCopier struct {
	deepCopy  bool
	skipField map[string]bool
}

type Option func(*deepCopier)

func WithDeepCopy(b bool) Option {
	return func(d *deepCopier) {
		d.deepCopy = b
	}
}

func WithSkipField(field string) Option {
	return func(d *deepCopier) {
		if d.skipField == nil {
			d.skipField = make(map[string]bool)
		}
		d.skipField[strings.ToLower(field)] = true
	}
}

func Copy(to any, from any, opts ...Option) error {
	d := &deepCopier{
		deepCopy: false,
	}
	for _, o := range opts {
		o(d)
	}

	destValue := reflect.ValueOf(to)
	if destValue.Kind() != reflect.Ptr {
		return errors.New("to not ptr type")
	} else {
		destValue = destValue.Elem()
		if destValue.Kind() != reflect.Struct {
			return errors.New("to not struct type")
		}
	}

	fromValue := reflect.Indirect(reflect.ValueOf(from))
	if fromValue.Kind() != reflect.Struct {
		return errors.New("from not struct type")
	}

	return d.copier(destValue, fromValue, "")
}

func (d *deepCopier) copier(to, from reflect.Value, fullName string) error {
	numField := to.NumField()
	for i := 0; i < numField; i++ {
		var (
			dstFieldValue = to.Field(i)
			dstFieldType  = to.Type().Field(i)
			srcFieldValue = reflect.Indirect(from.FieldByName(dstFieldType.Name))
		)
		if d.skipField != nil {
			if _, ok := d.skipField[getFullName(fullName, dstFieldType.Name)]; ok {
				continue
			}
		}

		if !srcFieldValue.IsValid() { // 找不到字段
			continue
		}
		if srcFieldValue.IsZero() {
			continue
		}

		if !dstFieldValue.CanSet() {
			fmt.Println("CanSet  ", dstFieldType.Name)
			continue
		}
		if !dstFieldType.IsExported() {
			continue
		}

		switch dstFieldValue.Kind() {
		case reflect.Struct:
			if err := d.copier(dstFieldValue, srcFieldValue, getFullName(fullName, dstFieldType.Name)); err != nil {
				return err
			}
		case reflect.Slice:
			if err := d.sliceCopy(dstFieldValue, srcFieldValue, fullName); err != nil {
				return err
			}
		case reflect.Ptr:
			if err := d.prtCopy(dstFieldValue, srcFieldValue, getFullName(fullName, dstFieldType.Name)); err != nil {
				return err
			}
		case reflect.Map:
			if err := d.mapCopy(dstFieldValue, srcFieldValue, fullName); err != nil {
				return err
			}
		default:
			if dstFieldType.Type.AssignableTo(srcFieldValue.Type()) {
				d.basicCopy(dstFieldValue, srcFieldValue, fullName)
				//dstFieldValue.Set(srcFieldValue)
			}
		}
	}
	return nil
}

func (d *deepCopier) mapCopy(dst, from reflect.Value, fullName string) error {
	if !d.deepCopy {
		dst.Set(from)
		return nil
	}

	makeMap := reflect.MakeMap(dst.Type())
	iter := from.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()

		makeMap.SetMapIndex(k, v)
	}
	dst.Set(makeMap)
	return nil
}

func (d *deepCopier) sliceCopy(dst, from reflect.Value, fullName string) error {
	if !d.deepCopy {
		dst.Set(from)
		return nil
	}

	makeSlice := reflect.MakeSlice(dst.Type(), from.Len(), from.Cap())
	reflect.Copy(makeSlice, from)
	dst.Set(makeSlice)
	return nil
}

func (d *deepCopier) prtCopy(dst, from reflect.Value, fullName string) error {
	if dst.IsNil() {
		dst.Set(reflect.New(dst.Type().Elem()))
	}
	dst = reflect.Indirect(dst)
	if !d.deepCopy {
		dst.Set(from)
		return nil
	}

	return d.copier(dst, from, fullName)
}

func (d *deepCopier) basicCopy(dst, from reflect.Value, fullName string) {
	switch dst.Kind() {
	case reflect.String:
		v := from.String()
		dst.SetString(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := from.Int()
		dst.SetInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := from.Uint()
		dst.SetUint(v)
	case reflect.Float32, reflect.Float64:
		v := from.Float()
		dst.SetFloat(v)
	case reflect.Bool:
		v := from.Bool()
		dst.SetBool(v)
	default:
		dst.Set(from)
	}
}

func getFullName(parent, child string) string {
	if len(parent) == 0 {
		return strings.ToLower(child)
	}

	return strings.ToLower(strings.Join([]string{parent, child}, "."))
}
