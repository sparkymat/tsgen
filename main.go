package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/samber/lo"
	"golang.org/x/mod/modfile"
)

func main() {
	fmt.Printf("args: %+v\n", os.Args)

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Printf("pwd: %s\n", pwd)

	var outBuffer bytes.Buffer

	cmd := exec.Command("go", "env", "GOMOD")
	cmd.Stdout = &outBuffer

	if err = cmd.Run(); err != nil {
		panic(err)
	}

	goModPath := outBuffer.String()

	goModPath = strings.Trim(goModPath, "\n")

	goModBytes, err := os.ReadFile(goModPath)
	if err != nil {
		panic(err)
	}

	goMod, err := modfile.Parse(goModPath, goModBytes, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("goMod: %s\n", goMod.Module.Mod.String())

	relPath, err := filepath.Rel(filepath.Dir(goModPath), pwd)
	if err != nil {
		panic(err)
	}

	fmt.Printf("relPath: %s\n", relPath)

	entries, err := os.ReadDir(pwd)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileContents, err := os.ReadFile(filepath.Join(pwd, entry.Name()))
		if err != nil {
			panic(err)
		}

		var tokenFileset token.FileSet

		fileNode, err := parser.ParseFile(&tokenFileset, entry.Name(), fileContents, parser.AllErrors)
		if err != nil {
			panic(err)
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

		if len(typeDeclarations) > 0 {
			td := typeDeclarations[0]

			typeDecl := td.(*ast.GenDecl)

			structSpecs := lo.Filter(typeDecl.Specs, func(ts ast.Spec, _ int) bool {
				v, ok := ts.(*ast.TypeSpec)
				if !ok {
					return false
				}

				if _, ok := v.Type.(*ast.StructType); ok {
					return true
				}

				return false
			})

			for _, spec := range structSpecs {
				typeSpec := spec.(*ast.TypeSpec)

				name := typeSpec.Name.Name

				typedType := typeSpec.Type.(*ast.StructType)

				handleType(name, typedType.Fields.List)
			}
		}
	}
}

func handleType(name string, fields []*ast.Field) {
	fieldStrings := lo.Map(fields, func(f *ast.Field, _ int) string {
		identType := f.Type.(*ast.Ident)
		return fmt.Sprintf("%s: %s", f.Names[0].Name, identType.Name)
	})

	fmt.Printf("%s has fields: %v\n", name, fieldStrings)
}
