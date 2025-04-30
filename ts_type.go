package tsgen

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/samber/lo"
)

var ErrInvalidType = errors.New("invalid type")

type TSType struct {
	name   string
	fields map[string]string
}

func (m TSType) Name() string {
	return m.name
}

func (m TSType) Fields() map[string]string {
	return m.fields
}

func (m TSType) RenderedFieldsForClass() string {
	v := ""

	for name, fType := range m.fields {
		v += fmt.Sprintf("  public %s: %s;\n\n", name, fType)
	}

	return v
}

func (m TSType) RenderedFieldsForInterface() string {
	v := ""

	for name, fType := range m.fields {
		v += fmt.Sprintf("  %s: %s;\n", name, fType)
	}

	return v
}

func isModel(fType string) bool {
	return !lo.Contains([]string{"string", "number", "boolean"}, fType)
}

func (m TSType) RenderedFieldAssignments() string {
	v := ""

	for name, fType := range m.fields {
		trimmedName := strings.TrimSuffix(name, "?")

		if strings.HasSuffix(name, "?") {
			v += fmt.Sprintf("    if (json.%s) {\n  ", trimmedName)
		}

		if isModel(fType) {
			v += fmt.Sprintf("    this.%s = new %s(json.%s);\n", trimmedName, fType, trimmedName)
		} else {
			v += fmt.Sprintf("    this.%s = json.%s;\n", trimmedName, trimmedName)
		}

		if strings.HasSuffix(name, "?") {
			v += "    }\n"
		}
	}

	return v
}

func (m TSType) Imports() string {
	v := ""

	models := lo.Uniq(
		lo.Filter(
			lo.Values(m.fields),
			func(fType string, _ int) bool { return isModel(fType) },
		),
	)

	for _, model := range models {
		v += fmt.Sprintf("import { %s } from './%s';\n", model, model)
	}

	if len(models) > 0 {
		v += "\n"
	}

	return v
}

func StructToTSType(v any, addID bool) (TSType, error) {
	val := reflect.ValueOf(v)

	if val.Type().Kind() != reflect.Struct {
		return TSType{}, ErrInvalidType
	}

	tt := TSType{
		name:   val.Type().Name(),
		fields: map[string]string{},
	}

	if addID {
		tt.fields["id"] = "string"
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

		tt.fields[fieldName] = fieldType
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
