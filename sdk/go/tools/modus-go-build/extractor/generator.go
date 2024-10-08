package extractor

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/types"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

type Field struct {
	Name          string
	Type          string
	Tag           string
	IsComplexType bool
	IsSliceType   bool
}

type TemplateData struct {
	StructName       string
	HasComplexFields bool
	Fields           []Field
	OmittedFuncs     []string
}

func isSlice(t types.Type) bool {
	_, ok := t.(*types.Slice)
	return ok
}

func hasNamedObject(typ types.Type) bool {
	switch t := typ.(type) {
	case *types.Slice:
		return hasNamedObject(t.Elem())
	case *types.Array:
		return hasNamedObject(t.Elem())
	case *types.Pointer:
		return hasNamedObject(t.Elem())
	case *types.Map:
		return hasNamedObject(t.Key()) || hasNamedObject(t.Elem())
	case *types.Named:
		// If we find a named type, return true
		return true
	case *types.Basic:
		return false
	default:
		// Handle other cases, e.g., interfaces, maps, etc.
		// If these types can contain other types, you should add appropriate checks
		return false
	}
}

func formatType(typ types.Type) string {
	switch t := typ.(type) {
	case *types.Slice:
		return "[]" + formatType(t.Elem())
	case *types.Array:
		return fmt.Sprintf("[%d]%s", t.Len(), formatType(t.Elem()))
	case *types.Pointer:
		return "*" + formatType(t.Elem())
	case *types.Map:
		return "map[" + formatType(t.Key()) + "]" + formatType(t.Elem())
	case *types.Named:
		// Extract just the type name, omit the package name
		return t.Obj().Name()
	case *types.Basic:
		return t.Name()
	default:
		// Handle other cases, e.g., interfaces, maps, etc.
		return t.String()
	}
}

func getExportedTypes(pkgs map[string]*packages.Package) []TemplateData {
	var structs []TemplateData
	for _, pkg := range pkgs {
		for _, file := range pkg.Syntax {
			for _, decl := range file.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok {
					if attrs, toGen := getHypGenTypeName(genDecl); toGen {
						for _, spec := range genDecl.Specs {
							typeSpec, ok := spec.(*ast.TypeSpec)
							if !ok {
								continue
							}
							// Lookup the struct type
							obj := pkg.Types.Scope().Lookup(typeSpec.Name.Name)
							if obj == nil {
								continue
							}
							st, ok := obj.Type().Underlying().(*types.Struct)
							if !ok {
								continue
							}
							var fields []Field
							var hasComplexFields bool
							// Iterate through the fields of the struct
							for i := 0; i < st.NumFields(); i++ {
								field := st.Field(i)
								fieldName := field.Name()
								fieldType := formatType(field.Type())
								fieldIsSlice := isSlice(field.Type())
								fieldIsComplex := hasNamedObject(field.Type())
								if fieldIsComplex {
									hasComplexFields = true
								}
								fields = append(fields, Field{Name: fieldName, Type: fieldType, Tag: strings.TrimSuffix(strings.TrimPrefix(st.Tag(i), "`"), "`"), IsComplexType: fieldIsComplex, IsSliceType: fieldIsSlice})
							}

							splittedAttrs := strings.Split(attrs, " ")
							var omitFuncs []string
							for _, attr := range splittedAttrs {
								if strings.HasPrefix(attr, "omit:") {
									omitFuncs = strings.Split(strings.TrimPrefix(attr, "omit:"), ",")
								}
							}
							structs = append(structs, TemplateData{
								StructName:       typeSpec.Name.Name,
								HasComplexFields: hasComplexFields,
								Fields:           fields,
								OmittedFuncs:     omitFuncs,
							})
						}
					}
				}
			}
		}
	}

	return structs
}

func cleanup(dir, outputFileName string) error {
	err := os.Remove(filepath.Join(dir, outputFileName))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func generateFunctions(sourceDir string, pkgs map[string]*packages.Package) {
	structs := getExportedTypes(pkgs)

	if len(structs) == 0 {
		return // No structs to generate functions for
	}

	var buf bytes.Buffer
	headerStr := `
	package main
	import (
		"fmt"
		"encoding/json"
		"github.com/hypermodeinc/modus/sdk/go/pkg/dgraph"
	)

	const hostName = "dgraph"


	`

	buf.WriteString(headerStr)

	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}

	t1 := template.Must(template.New("upsert").Funcs(funcMap).Parse(upsertTmplStruct))
	t2 := template.Must(template.New("delete").Funcs(funcMap).Parse(deleteTmplStruct))
	t3 := template.Must(template.New("get").Funcs(funcMap).Parse(getTmplStruct))
	t4 := template.Must(template.New("query").Funcs(funcMap).Parse(queryTmplStruct))

	for _, data := range structs {
		var err error
		if !slices.Contains(data.OmittedFuncs, "upsert") {
			err = t1.Execute(&buf, data)
			if err != nil {
				fmt.Println("Error executing upsert template:", err)
				os.Exit(1)
			}
		}

		if !slices.Contains(data.OmittedFuncs, "delete") {
			err = t2.Execute(&buf, data)
			if err != nil {
				fmt.Println("Error executing delete template:", err)
				os.Exit(1)
			}
		}

		if !slices.Contains(data.OmittedFuncs, "get") {
			err = t3.Execute(&buf, data)
			if err != nil {
				fmt.Println("Error executing get template:", err)
				os.Exit(1)
			}
		}

		if !slices.Contains(data.OmittedFuncs, "query") {
			err = t4.Execute(&buf, data)
			if err != nil {
				fmt.Println("Error executing query template:", err)
				os.Exit(1)
			}
		}
	}

	formatted, err := format.Source(buf.Bytes())

	if err != nil {
		fmt.Println("Error formatting source:", err)
		fmt.Println("Unformatted source:\n", buf.String())
		os.Exit(1)
	}

	outputFileName := "hyp_functions_generated.go"
	err = cleanup(sourceDir, outputFileName)
	if err != nil {
		fmt.Println("Error cleaning up:", err)
		os.Exit(1)
	}

	err = os.WriteFile(filepath.Join(sourceDir, outputFileName), formatted, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		os.Exit(1)
	}

	fmt.Println("Methods generated in", outputFileName)
}
