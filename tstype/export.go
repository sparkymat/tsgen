package tstype

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

func (m TSType) RenderedFieldsForClass() string {
	v := ""

	for _, name := range m.orderedFieldNames {
		fType := m.fields[name]
		v += fmt.Sprintf("  public %s: %s;\n\n", name, fType)
	}

	return v
}

func (m TSType) RenderedFieldsForInterface() string {
	v := ""

	for _, name := range m.orderedFieldNames {
		fType := m.fields[name]
		v += fmt.Sprintf("  %s: %s;\n", name, fType)
	}

	return v
}

func (m TSType) RenderedFieldAssignments() string {
	v := ""

	for _, name := range m.orderedFieldNames {
		fType := m.fields[name]
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
