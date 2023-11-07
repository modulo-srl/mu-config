package generator

import (
	"go/constant"
	"go/types"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/modulo-srl/mu-config/go2cfg/distiller"
)

func TestGenerator(t *testing.T) {
	var tests = []struct {
		pkgDir   string
		typeName string
		filename string
		mode     DocTypesMode
	}{
		{"../testdata", "Embedding", "../testdata/embedding.jsonc", NoFields},
		{"../testdata", "Empty", "../testdata/empty.jsonc", NoFields},
		{"../testdata", "Nesting", "../testdata/nesting.jsonc", NoFields},
		{"../testdata", "Simple", "../testdata/simple.jsonc", NoFields},
		{"../testdata/multipkg", "MultiPackage", "../testdata/multipkg/multi_package.jsonc", NoFields},

		{"../testdata", "Embedding", "../testdata/embedding_basic_fields.jsonc", BasicFields},
		{"../testdata", "Nesting", "../testdata/nesting_basic_fields.jsonc", BasicFields},
		{"../testdata", "Simple", "../testdata/simple_basic_fields.jsonc", BasicFields},
		{"../testdata/multipkg", "MultiPackage", "../testdata/multipkg/multi_package_basic_fields.jsonc", BasicFields},

		{"../testdata", "Embedding", "../testdata/embedding_all_fields.jsonc", AllFields},
		{"../testdata", "Nesting", "../testdata/nesting_all_fields.jsonc", AllFields},
		{"../testdata", "Simple", "../testdata/simple_all_fields.jsonc", AllFields},
		{"../testdata/multipkg", "MultiPackage", "../testdata/multipkg/multi_package_all_fields.jsonc", AllFields},
	}

	whitespacesReplacer := strings.NewReplacer(" ", "◦", "\t", "———➞")
	for _, test := range tests {
		jsonc, err := Generate(test.pkgDir, test.typeName, test.mode)
		if err != nil {
			t.Fatal(err)
		}

		content, err := os.ReadFile(test.filename)
		if err != nil {
			t.Fatal(err)
		}

		want := string(content)

		if jsonc != want {
			t.Fatalf("Generated JSONC mismatch for %s struct:\n%s\n\nwant %s:\n%s",
				test.typeName,
				whitespacesReplacer.Replace(jsonc),
				test.filename,
				whitespacesReplacer.Replace(want))
		}
	}

	_, err := Generate("../testdata/invalid-path", "", AllFields)
	if err == nil {
		t.Fatalf("Generating for invalid path: expected error, got nil.")
	}

	_, err = Generate("../testdata", "invalid-struct", AllFields)
	if err == nil {
		t.Fatalf("Generating for invalid struct: expected error, got nil.")
	}
}

func TestGenerator_typeZero(t *testing.T) {
	tests := []struct {
		info *distiller.FieldInfo
		want interface{}
	}{
		{info: &distiller.FieldInfo{Type: nil, Layout: distiller.LayoutArray}, want: make([]interface{}, 0)},
		{info: &distiller.FieldInfo{Type: nil, Layout: distiller.LayoutMap}, want: make(map[interface{}]interface{})},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Bool], Layout: distiller.LayoutSingle}, want: false},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Int], Layout: distiller.LayoutSingle}, want: 0},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Int8], Layout: distiller.LayoutSingle}, want: int8(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Int16], Layout: distiller.LayoutSingle}, want: int16(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Int32], Layout: distiller.LayoutSingle}, want: int32(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Int64], Layout: distiller.LayoutSingle}, want: int64(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Uint], Layout: distiller.LayoutSingle}, want: uint(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Uint8], Layout: distiller.LayoutSingle}, want: uint8(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Uint16], Layout: distiller.LayoutSingle}, want: uint16(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Uint32], Layout: distiller.LayoutSingle}, want: uint32(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Uint64], Layout: distiller.LayoutSingle}, want: uint64(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Uintptr], Layout: distiller.LayoutSingle}, want: uintptr(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Float32], Layout: distiller.LayoutSingle}, want: float32(0.0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Float64], Layout: distiller.LayoutSingle}, want: 0.0},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Complex64], Layout: distiller.LayoutSingle}, want: complex64(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Complex128], Layout: distiller.LayoutSingle}, want: complex128(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.String], Layout: distiller.LayoutSingle}, want: constant.MakeString("")},
	}

	for _, test := range tests {
		zero := typeZero(test.info)

		if !reflect.DeepEqual(zero, test.want) {
			t.Fatalf("Zero value mismatch for type %v: got %v, want %v",
				test.info.Type.String(), zero, test.want)
		}
	}
}
