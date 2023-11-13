package distiller

import (
	"go/ast"
	"golang.org/x/tools/go/packages"
	"testing"
)

func LoadPackage(t *testing.T, patterns ...string) []*packages.Package {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles |
			packages.NeedImports | packages.NeedDeps | packages.NeedSyntax |
			packages.NeedTypes | packages.NeedTypesInfo,
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		t.Fatal(err)
	}

	if len(pkgs) == 0 {
		t.Fatalf("Expected at least one package reading: %v", patterns)
	}

	return pkgs
}

func GetFieldsInfo(t *testing.T, patterns []string) []*FieldInfo {
	pkgs := LoadPackage(t, patterns...)
	var fields []*FieldInfo

	for _, pattern := range patterns {
		_, err := NewPackageInfo(pattern, "")
		if err != nil {
			t.Fatalf("Error loading package %s: %v", pattern, err)
		}
	}

	for _, pkg := range pkgs {
		for _, astFile := range pkg.Syntax {
			for _, decl := range astFile.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}

				ast.Inspect(genDecl, func(node ast.Node) bool {
					var field *ast.Field
					field, ok = node.(*ast.Field)
					if !ok {
						return true
					}

					fields = append(fields, NewFieldInfo(field, pkg))

					return true
				})
			}
		}
	}
	return fields
}
