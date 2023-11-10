package renderers

import (
	"github.com/modulo-srl/mu-config/go2cfg/distiller"
	"github.com/modulo-srl/mu-config/go2cfg/testutils"
	"go/constant"
	"go/types"
	"reflect"
	"strings"
	"testing"
)

func TestFieldInfo_renderDoc(t *testing.T) {
	info := testutils.GetFieldsInfo(t, []string{"../testdata"})
	testTable := []struct {
		types bool
		want  string
	}{
		{types: true, want: "// int - Identifier documentation block.\n"},
		{types: true, want: "// bool - Enabled comment line.\n"},
		{types: true, want: "// uint32\n"},
		{types: true, want: "// testdata.Embedded - Embedded documentation block.\n"},
		{types: true, want: "// float32 - Position comment line.\n"},
		{types: true, want: "// float32 - Velocity documentation block.\n"},
		{types: true, want: "// float32\n"},
		{types: true, want: "// string - Shadowing field.\n"},
		{types: true, want: "// string - Field A\n"},
		{types: true, want: "// int - Field B\n"},
		{types: true, want: "// testdata.EmptySubType\n"},
		{types: true, want: "// []testdata.EmptySubType\n"},
		{types: true, want: "// string - Name describes the protocol name.\n// Multiple line documentation test.\n// Protocol name.\n"},
		{types: true, want: "// int - Major version.\n"},
		{types: true, want: "// int - Minor version.\n"},
		{types: true, want: "// string - Remote IP address.\n"},
		{types: true, want: "// int - Remote port.\n"},
		{types: true, want: "// testdata.Protocol - Default protocol.\n"},
		{types: true, want: "// []testdata.Protocol - Optional supported protocols.\n"},
		{types: true, want: "// string - Name of the user documentation block.\n// User name comment.\n"},
		{types: true, want: "// string - User surname comment.\n"},
		{types: true, want: "// int - Age documentation block.\n// User age.\n"},
		{types: true, want: "// int - Number of stars achieved.\n"},
		{types: true, want: "// []string - Addresses comment.\n"},
		{types: true, want: "// map[string]string - User tags.\n"},
		{
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

		{types: false, want: "// Identifier documentation block.\n"},
		{types: false, want: "// Enabled comment line.\n"},
		{types: false, want: ""},
		{types: false, want: "// Embedded documentation block.\n"},
		{types: false, want: "// Position comment line.\n"},
		{types: false, want: "// Velocity documentation block.\n"},
		{types: false, want: ""},
		{types: false, want: "// Shadowing field.\n"},
		{types: false, want: "// Field A\n"},
		{types: false, want: "// Field B\n"},
		{types: false, want: ""},
		{types: false, want: ""},
		{types: false, want: "// Name describes the protocol name.\n// Multiple line documentation test.\n// Protocol name.\n"},
		{types: false, want: "// Major version.\n"},
		{types: false, want: "// Minor version.\n"},
		{types: false, want: "// Remote IP address.\n"},
		{types: false, want: "// Remote port.\n"},
		{types: false, want: "// Default protocol.\n"},
		{types: false, want: "// Optional supported protocols.\n"},
		{types: false, want: "// Name of the user documentation block.\n// User name comment.\n"},
		{types: false, want: "// User surname comment.\n"},
		{types: false, want: "// Age documentation block.\n// User age.\n"},
		{types: false, want: "// Number of stars achieved.\n"},
		{types: false, want: "// Addresses comment.\n"},
		{types: false, want: "// User tags.\n"},
		{
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
	}

	whitespacesReplacer := strings.NewReplacer(" ", "◦", "\t", "———➞", "\n", "⏎\n")
	for i, test := range testTable {
		doc := renderDoc(info[i%len(info)], "", "//", test.types)
		if doc != test.want {
			t.Fatalf("FormatDoc return mismatch:\ngot:%v\nwant:%v\n",
				whitespacesReplacer.Replace(doc),
				whitespacesReplacer.Replace(test.want))
		}
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
