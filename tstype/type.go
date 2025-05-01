package tstype

import (
	"errors"
	"reflect"

	"github.com/samber/lo"
)

var ErrInvalidType = errors.New("invalid type")

type TSType struct {
	name              string
	fields            map[string]string
	orderedFieldNames []string
}

func New(name string) TSType {
	return TSType{
		name: name,
	}
}

func (t *TSType) Name() string {
	return t.name
}

func (t *TSType) Fields() map[string]string {
	return t.fields
}

func (t *TSType) AddField(name string, typeName string) {
	t.fields[name] = typeName
	t.orderedFieldNames = append(t.orderedFieldNames, name)
}

func isModel(fType string) bool {
	return !lo.Contains([]string{"string", "number", "boolean"}, fType)
}

func StructToTSType(v any, addID bool) (TSType, error) {
	val := reflect.ValueOf(v)

	if val.Type().Kind() != reflect.Struct {
		return TSType{}, ErrInvalidType
	}

	tt := TSType{
		name:              val.Type().Name(),
		fields:            map[string]string{},
		orderedFieldNames: []string{},
	}

	if addID {
		tt.AddField("id", "string")
	}

	for fi := range val.Type().NumField() {
		f := val.Type().Field(fi)

		fieldName := f.Tag.Get("json")

		if fieldName == "" {
			fieldName = f.Tag.Get("query")
		}

		if fieldName == "" {
			fieldName = f.Tag.Get("form")
		}

		if fieldName == "" {
			fieldName = f.Name
		}

		fieldName, fieldType := fieldToTSType(fieldName, f.Type)

		tt.AddField(fieldName, fieldType)
	}

	return tt, nil
}

func fieldToTSType(name string, goType reflect.Type) (string, string) {
	switch goType.Kind() {
	case reflect.String:
		return name, "string"
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex128,
		reflect.Complex64:
		return name, "number"
	case reflect.Bool:
		return name, "boolean"
	case reflect.Pointer:
		return fieldToTSType(name+"?", goType.Elem())
	case reflect.Array, reflect.Slice:
		updatedName, elemType := fieldToTSType(name, goType.Elem())

		return updatedName, elemType + "[]"
	case reflect.Struct:
		return name, goType.Name()
	case reflect.Map,
		reflect.Invalid,
		reflect.Func,
		reflect.Interface,
		reflect.Chan,
		reflect.Uintptr,
		reflect.UnsafePointer:
		fallthrough
	default:
		return name, "any"
	}
}
