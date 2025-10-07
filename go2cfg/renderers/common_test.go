package renderers

import (
	"go/constant"
	"go/types"
	"reflect"
	"strings"
	"testing"

	"github.com/modulo-srl/mu-config/go2cfg/distiller"
)

func TestFieldInfo_renderDoc(t *testing.T) {
	info := distiller.GetFieldsInfo(t, []string{"../testdata"})
	testTable := []struct {
		name  string
		types bool
		want  string
	}{
		{name: "type-int-block", types: true, want: "// int - Identifier documentation block.\n"},
		{name: "type-bool-line", types: true, want: "// bool - Enabled comment line.\n"},
		{name: "type-uint32-none", types: true, want: "// uint32\n"},
		{name: "type-embedded-block", types: true, want: "// testdata.Embedded - Embedded documentation block.\n"},
		{name: "type-float32-line", types: true, want: "// float32 - Position comment line.\n"},
		{name: "type-float32-block", types: true, want: "// float32 - Velocity documentation block.\n"},
		{name: "type-float32-none", types: true, want: "// float32\n"},
		{name: "type-string-shadowing", types: true, want: "// string - Shadowing field.\n"},
		{name: "type-string-line", types: true, want: "// string - Field A\n"},
		{name: "type-int-line", types: true, want: "// int - Field B\n"},
		{name: "type-subtype-none", types: true, want: "// testdata.EmptySubType\n"},
		{name: "type-subtype-slice-none", types: true, want: "// []testdata.EmptySubType\n"},
		{name: "type-string-multi-line", types: true, want: "// string - Name describes the protocol name.\n// Multiple line documentation test.\n// Protocol name.\n"},
		{name: "type-subtype-int-line", types: true, want: "// int - Major version.\n"},
		{name: "type-subtype-int-line2", types: true, want: "// int - Minor version.\n"},
		{name: "type-nesting-string-block", types: true, want: "// string - Remote IP address.\n"},
		{name: "type-nesting-int-block", types: true, want: "// int - Remote port.\n"},
		{name: "type-nesting-subtype-line", types: true, want: "// testdata.Protocol - Default protocol.\n"},
		{name: "type-nesting-subtype-slice-line", types: true, want: "// []testdata.Protocol - Optional supported protocols.\n"},
		{name: "type-string-block-line", types: true, want: "// string - Name of the user documentation block.\n// User name comment.\n"},
		{name: "type-string-block", types: true, want: "// string - User surname comment.\n"},
		{name: "type-int-block-line", types: true, want: "// int - Age documentation block.\n// User age.\n"},
		{name: "type-int-json-line", types: true, want: "// int - Number of stars achieved.\n"},
		{name: "type-string-array-block", types: true, want: "// []string - Addresses comment.\n"},
		{name: "type-map-line", types: true, want: "// map[string]string - User tags.\n"},
		{
			name:  "type-const-multi-line",
			types: true,
			want: `// testdata.ConstType - Type documentation block.
// Type of constant.
// Allowed values:
// ConstTypeA =   0  ConstTypeA doc block. ConstTypeA comment.
// ConstTypeB =   1  ConstTypeB comment.
// ConstTypeC =   2  ConstTypeC doc block. ConstTypeC comment.
// ConstTypeD =  32  ConstTypeD doc block.
// ConstTypeE =  64  ConstTypeE doc block. ConstTypeE comment.
// ConstTypeF = 128  ConstTypeF doc block. ConstTypeF comment.
`,
		},
		{name: "type-multi-vars-x", types: true, want: "// float64 - X, Y documentation block.\n// Coordinates.\n"},
		{name: "type-multi-vars-y", types: true, want: "// float64 - X, Y documentation block.\n// Coordinates.\n"},

		{name: "int-block", types: false, want: "// Identifier documentation block.\n"},
		{name: "bool-line", types: false, want: "// Enabled comment line.\n"},
		{name: "uint32-none", types: false, want: ""},
		{name: "embedded-block", types: false, want: "// Embedded documentation block.\n"},
		{name: "float32-line", types: false, want: "// Position comment line.\n"},
		{name: "float32-block", types: false, want: "// Velocity documentation block.\n"},
		{name: "float32-none", types: false, want: ""},
		{name: "string-shadowing", types: false, want: "// Shadowing field.\n"},
		{name: "string-line", types: false, want: "// Field A\n"},
		{name: "int-line", types: false, want: "// Field B\n"},
		{name: "subtype-none", types: false, want: ""},
		{name: "subtype-slice-none", types: false, want: ""},
		{name: "string-multi-line", types: false, want: "// Name describes the protocol name.\n// Multiple line documentation test.\n// Protocol name.\n"},
		{name: "subtype-int-line", types: false, want: "// Major version.\n"},
		{name: "subtype-int-line2", types: false, want: "// Minor version.\n"},
		{name: "nesting-string-block", types: false, want: "// Remote IP address.\n"},
		{name: "nesting-int-block", types: false, want: "// Remote port.\n"},
		{name: "nesting-subtype-line", types: false, want: "// Default protocol.\n"},
		{name: "nesting-subtype-slice", types: false, want: "// Optional supported protocols.\n"},
		{name: "string-block-line", types: false, want: "// Name of the user documentation block.\n// User name comment.\n"},
		{name: "string-block", types: false, want: "// User surname comment.\n"},
		{name: "int-block-line", types: false, want: "// Age documentation block.\n// User age.\n"},
		{name: "int-json-line", types: false, want: "// Number of stars achieved.\n"},
		{name: "string-array-block", types: false, want: "// Addresses comment.\n"},
		{name: "map-line", types: false, want: "// User tags.\n"},
		{
			name:  "const-multi-line",
			types: false,
			want: `// Type documentation block.
// Type of constant.
// Allowed values:
// ConstTypeA =   0  ConstTypeA doc block. ConstTypeA comment.
// ConstTypeB =   1  ConstTypeB comment.
// ConstTypeC =   2  ConstTypeC doc block. ConstTypeC comment.
// ConstTypeD =  32  ConstTypeD doc block.
// ConstTypeE =  64  ConstTypeE doc block. ConstTypeE comment.
// ConstTypeF = 128  ConstTypeF doc block. ConstTypeF comment.
`,
		},
		{name: "multi-vars-x", types: false, want: "// X, Y documentation block.\n// Coordinates.\n"},
		{name: "multi-vars-y", types: false, want: "// X, Y documentation block.\n// Coordinates.\n"},
	}

	whitespacesReplacer := strings.NewReplacer(" ", "◦", "\t", "———➞", "\n", "⏎\n")
	for i, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			doc := renderDoc(info[i%len(info)], "", "//", test.types)
			if doc != test.want {
				t.Fatalf("FormatDoc return mismatch:\ngot:%v\nwant:%v\n",
					whitespacesReplacer.Replace(doc),
					whitespacesReplacer.Replace(test.want))
			}
		})
	}
}

func TestRenderers_lastIndexOf(t *testing.T) {
	tests := []struct {
		slice []string
		value string
		want  int
	}{
		{slice: []string{"a", "b", "c", "b"}, value: "a", want: 0},
		{slice: []string{"a", "b", "c", "b"}, value: "b", want: 3},
		{slice: []string{"a", "b", "c", "b"}, value: "c", want: 2},
		{slice: []string{"a", "b", "c", "b"}, value: "d", want: -1},
	}

	for _, test := range tests {
		got := lastIndexOf(test.slice, test.value)
		if got != test.want {
			t.Fatalf("index mismatch: got %v, want %v", got, test.want)
		}
	}
}

func TestRenderers_typeZero(t *testing.T) {
	tests := []struct {
		name string
		info *distiller.FieldInfo
		want interface{}
	}{
		{name: "NilArray", info: &distiller.FieldInfo{Type: nil, Layout: distiller.LayoutArray}, want: make([]interface{}, 0)},
		{name: "NilMap", info: &distiller.FieldInfo{Type: nil, Layout: distiller.LayoutMap}, want: make(map[interface{}]interface{})},
		{name: "Bool", info: &distiller.FieldInfo{Type: types.Typ[types.Bool], Layout: distiller.LayoutSingle}, want: false},
		{name: "Int", info: &distiller.FieldInfo{Type: types.Typ[types.Int], Layout: distiller.LayoutSingle}, want: 0},
		{name: "Int8", info: &distiller.FieldInfo{Type: types.Typ[types.Int8], Layout: distiller.LayoutSingle}, want: int8(0)},
		{name: "Int16", info: &distiller.FieldInfo{Type: types.Typ[types.Int16], Layout: distiller.LayoutSingle}, want: int16(0)},
		{name: "Int32", info: &distiller.FieldInfo{Type: types.Typ[types.Int32], Layout: distiller.LayoutSingle}, want: int32(0)},
		{name: "Int64", info: &distiller.FieldInfo{Type: types.Typ[types.Int64], Layout: distiller.LayoutSingle}, want: int64(0)},
		{name: "Uint", info: &distiller.FieldInfo{Type: types.Typ[types.Uint], Layout: distiller.LayoutSingle}, want: uint(0)},
		{name: "Uint8", info: &distiller.FieldInfo{Type: types.Typ[types.Uint8], Layout: distiller.LayoutSingle}, want: uint8(0)},
		{name: "Uint16", info: &distiller.FieldInfo{Type: types.Typ[types.Uint16], Layout: distiller.LayoutSingle}, want: uint16(0)},
		{name: "Uint32", info: &distiller.FieldInfo{Type: types.Typ[types.Uint32], Layout: distiller.LayoutSingle}, want: uint32(0)},
		{name: "Uint64", info: &distiller.FieldInfo{Type: types.Typ[types.Uint64], Layout: distiller.LayoutSingle}, want: uint64(0)},
		{name: "Uintptr", info: &distiller.FieldInfo{Type: types.Typ[types.Uintptr], Layout: distiller.LayoutSingle}, want: uintptr(0)},
		{name: "Float32", info: &distiller.FieldInfo{Type: types.Typ[types.Float32], Layout: distiller.LayoutSingle}, want: float32(0.0)},
		{name: "Float64", info: &distiller.FieldInfo{Type: types.Typ[types.Float64], Layout: distiller.LayoutSingle}, want: 0.0},
		{name: "Complex64", info: &distiller.FieldInfo{Type: types.Typ[types.Complex64], Layout: distiller.LayoutSingle}, want: complex64(0)},
		{name: "Complex128", info: &distiller.FieldInfo{Type: types.Typ[types.Complex128], Layout: distiller.LayoutSingle}, want: complex128(0)},
		{name: "String", info: &distiller.FieldInfo{Type: types.Typ[types.String], Layout: distiller.LayoutSingle}, want: constant.MakeString("")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			zero := typeZero(test.info)

			if !reflect.DeepEqual(zero, test.want) {
				t.Fatalf("Zero value mismatch for type %v: got %v, want %v",
					test.info.Type.String(), zero, test.want)
			}
		})
	}
}
