package extractor

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"slices"
	"strings"

	"github.com/samber/lo"
)

func ExtractStructTypes(fileName string, fileContents string) ([]StructType, error) {

	var tokenFileSet token.FileSet

	fileNode, err := parser.ParseFile(&tokenFileSet, fileName, fileContents, parser.DeclarationErrors)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file '%s': %w", fileName, err)
	}

	typeDeclarations := lo.Filter(fileNode.Decls, func(d ast.Decl, _ int) bool {
		switch v := d.(type) {
		case *ast.GenDecl:
			if v.Tok == token.TYPE {
				return true
			}

			return false
		default:
			return false
		}
	})

	structTypes := []StructType{}

	for _, typeDecl := range typeDeclarations {
		genDecl, ok := typeDecl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			typeName := typeSpec.Name.Name

			structTypeSpec, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			structType := structTypeFromSpec(typeName, structTypeSpec)

			structTypes = append(structTypes, structType)
		}
	}

	return structTypes, nil
}

func structTypeFromSpec(name string, spec *ast.StructType) StructType {
	t := StructType{}

	t.Name = name

	fields := []StructField{}

	for _, f := range spec.Fields.List {
		structField := structFieldFromSpec(f)

		fields = append(fields, structField)
	}

	t.Fields = fields

	return t
}

func structFieldFromSpec(f *ast.Field) StructField {
	t := StructField{}

	typeIdent, ok := f.Type.(*ast.Ident)
	if ok {
		switch {
		case typeIdent.Name == "string":
			t.Type = TypeString
		case slices.Contains([]string{
			"uint8",
			"uint16",
			"uint32",
			"uint64",
			"int8",
			"int16",
			"int32",
			"int64",
			"float32",
			"float64",
			"complex64",
			"complex128",
		}, typeIdent.Name):
			t.Type = TypeNumber
		case typeIdent.Name == "bool":
			t.Type = TypeBoolean
		case typeIdent.Obj != nil:
			t.Type = TypeStruct
		}
	}

	t.Name = f.Names[0].Name

	if f.Tag != nil {
		tagValue := strings.Trim(f.Tag.Value, "`")
		t.Tags = strings.Split(tagValue, " ")

		for _, tagString := range t.Tags {
			words := strings.Split(tagString, ":")

			if t.TagMap == nil {
				t.TagMap = map[string][]string{}
			}

			cleanedValuesString := strings.Trim(words[1], "\"")
			values := strings.Split(cleanedValuesString, ",")

			t.TagMap[words[0]] = values
		}
	}

	return t
}
