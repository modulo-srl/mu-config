package distiller

import (
	"fmt"
	"go/ast"
	"go/constant"
	"reflect"
	"strings"
	"testing"

	"github.com/modulo-srl/mu-config/go2cfg/ordered"
	"github.com/modulo-srl/mu-config/go2cfg/testdata/multipkg/network"
)

type StructInfoMatch struct {
	Package     string
	Name        string
	FieldsCount int
	Doc         string
	Defaults    map[string]interface{}
}

func (s *StructInfoMatch) String() string {
	return fmt.Sprintf("Package: %s\nName: \"%s\"\nFieldCount: %v\nDoc: \"%v\"\nDefaults: %+v",
		s.Package, s.Name, s.FieldsCount, strings.ReplaceAll(s.Doc, "\n", "\\n"), s.Defaults)
}

func TestStructInfo(t *testing.T) {
	testStructInfo(t, "../testdata", []*StructInfoMatch{
		// testdata/embedding.go
		{
			Package:     "github.com/modulo-srl/mu-config/go2cfg/testdata",
			Name:        "Embedded",
			Doc:         "Embedded test struct.\n",
			FieldsCount: 3,
			Defaults:    nil,
		},
		{
			Package:     "github.com/modulo-srl/mu-config/go2cfg/testdata",
			Name:        "Embedding",
			Doc:         "Embedding test struct.\n",
			FieldsCount: 5,
			Defaults:    nil,
		},
		// testdata/empty.go
		{
			Package:     "github.com/modulo-srl/mu-config/go2cfg/testdata",
			Name:        "Empty",
			Doc:         "Empty empty test struct.\n",
			FieldsCount: 0,
			Defaults:    nil,
		},
		// testdata/empty_defs.go
		{
			Package:     "github.com/modulo-srl/mu-config/go2cfg/testdata",
			Name:        "EmptySubType",
			Doc:         "EmptySubType define a struct with non-initialized fields.\n",
			FieldsCount: 2,
			Defaults:    nil,
		},
		{
			Package:     "github.com/modulo-srl/mu-config/go2cfg/testdata",
			Name:        "EmptyDefs",
			Doc:         "EmptyDefs define a struct with non-initialized fields.\n",
			FieldsCount: 2,
			Defaults:    nil,
		},
		// testdata/nesting.go
		{
			Package:     "github.com/modulo-srl/mu-config/go2cfg/testdata",
			Name:        "Protocol",
			Doc:         "Protocol defines a network protocol and version.\n",
			FieldsCount: 3,
			Defaults:    nil,
		},
		{
			Package:     "github.com/modulo-srl/mu-config/go2cfg/testdata",
			Name:        "Nesting",
			Doc:         "Nesting checks for correct struct nesting.\n",
			FieldsCount: 4,
			Defaults:    nil,
		},
		// testdata/simple.go
		{
			Package:     "github.com/modulo-srl/mu-config/go2cfg/testdata",
			Name:        "Simple",
			Doc:         "Simple defines a simple user.\n",
			FieldsCount: 9,
			Defaults:    nil,
		},
	})
}

func TestStructInfoMultiPackage(t *testing.T) {
	testStructInfo(t, "../testdata/multipkg", []*StructInfoMatch{
		// testdata/multipkg/multi_package.go
		{
			Package:     "github.com/modulo-srl/mu-config/go2cfg/testdata/multipkg",
			Name:        "MultiPackage",
			Doc:         "MultiPackage tests the multi-package and import aliasing case.\n",
			FieldsCount: 2,
			Defaults:    nil,
		},
	})
}

func TestStructInfoDefaults(t *testing.T) {
	tags := ordered.NewMap()
	tags.Append(constant.MakeString("Key1").String(), constant.MakeString("Value1"))
	tags.Append(constant.MakeString("Key2").String(), constant.MakeString("Value2"))
	tags.Append(constant.MakeString("Key3").String(), constant.MakeString("Value3"))

	testStructInfoDefaults(t, "../testdata", "Simple", map[string]interface{}{
		"Name":       constant.MakeString("John"),
		"Age":        constant.MakeInt64(30),
		"StarsCount": constant.MakeInt64(5),
		"Addresses": []interface{}{
			constant.MakeString("Address 1"),
			constant.MakeString("Address 2"),
			constant.MakeString("Address 3"),
		},
		"Tags": tags,
		"X":    constant.MakeFloat64(1.0),
		"Y":    constant.MakeFloat64(2.0),
	})
}

func TestStructInfoDefaultsMultiPackage(t *testing.T) {
	testStructInfoDefaults(t, "../testdata/multipkg", "MultiPackage", map[string]interface{}{
		"NetStatus": map[string]interface{}{
			"Connected": constant.MakeBool(true),
			"State":     constant.MakeInt64(int64(network.StateDisconnected)),
		},
		"Info": map[string]interface{}{
			"PacketLoss":    constant.MakeInt64(64),
			"RoundTripTime": constant.MakeInt64(123),
		},
	})
}

func testStructInfo(t *testing.T, pattern string, want []*StructInfoMatch) {
	pkgs := LoadPackage(t, pattern)

	var structs []*StructInfo

	for _, pkg := range pkgs {
		for _, astFile := range pkg.Syntax {
			for _, decl := range astFile.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}

				ast.Inspect(genDecl, func(node ast.Node) bool {
					var typeSpec *ast.TypeSpec
					typeSpec, ok = node.(*ast.TypeSpec)
					if !ok {
						return true
					}

					if _, ok = typeSpec.Type.(*ast.StructType); !ok {
						return true
					}

					structs = append(structs, NewStructInfo(genDecl, pkg))
					return true
				})
			}
		}
	}

	if len(structs) != len(want) {
		t.Fatalf("Parsed %d structs, want %d.", len(structs), len(want))
	}

	for i, s := range structs {
		if s.Package.PkgPath != want[i].Package ||
			s.Name != want[i].Name || s.Doc != want[i].Doc ||
			len(s.Fields) != want[i].FieldsCount ||
			!reflect.DeepEqual(s.Defaults, want[i].Defaults) {
			t.Fatalf("Parsed struct mismatch:\n%s\n\nwant:\n%s\n", s, want[i])
		}
	}
}

func testStructInfoDefaults(t *testing.T, pattern string, name string, want map[string]interface{}) {
	pkgs := LoadPackage(t, pattern)

	for _, pkg := range pkgs {
		for _, astFile := range pkg.Syntax {
			for _, decl := range astFile.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}

				ast.Inspect(genDecl, func(node ast.Node) bool {
					var typeSpec *ast.TypeSpec
					typeSpec, ok = node.(*ast.TypeSpec)
					if !ok {
						return true
					}

					if _, ok = typeSpec.Type.(*ast.StructType); !ok {
						return true
					}

					if typeSpec.Name.Name != name {
						return true
					}

					s := NewStructInfo(genDecl, pkg)
					if err := s.ParseDefaultsMethod(); err != nil {
						t.Fatal(err)
					}

					if !reflect.DeepEqual(s.Defaults, want) {
						t.Fatalf("Struct %s defaults mismatch:\n%+v\n\nwant:\n%+v", s.Name, s.Defaults, want)
					}

					return false
				})
			}
		}
	}
}
