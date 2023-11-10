package distiller

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/modulo-srl/mu-config/go2cfg/testutils"
)

type FieldInfoMatch struct {
	Type       string
	Name       string
	Layout     FieldLayout
	IsEmbedded bool
	Tags       map[string]string
	Doc        string
}

func (f *FieldInfoMatch) String() string {
	return fmt.Sprintf("Type: %s\nName: \"%s\"\nLayout: %v\nIsEmbedded: %v\nTags: %+v\nDoc: \"%v\"\n",
		f.Type, f.Name, f.Layout, f.IsEmbedded, f.Tags, strings.ReplaceAll(f.Doc, "\n", "\\n"))
}

func TestFieldInfo_String(t *testing.T) {
	info := testutils.GetFieldsInfo(t, []string{"../testdata"})
	want := []string{
		`Type: int
Name: "Identifier"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[json:id]
Doc: "Identifier documentation block.\n"
`,
		`Type: bool
Name: "Enabled"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: "Enabled comment line.\n"
`,
		`Type: uint32
Name: "Reserved"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[json:reserved]
Doc: ""
`,
		`Type: github.com/modulo-srl/mu-config/go2cfg/testdata.Embedded
Name: ""
Layout: 0
Element type: <nil>
IsEmbedded: true
Tags: map[]
Doc: "Embedded documentation block.\n"
`,
		`Type: float32
Name: "Position"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[json:position]
Doc: "Position comment line.\n"
`,
		`Type: float32
Name: "Velocity"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[json:velocity]
Doc: "Velocity documentation block.\n"
`,
		`Type: float32
Name: "Acceleration"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[json:accel]
Doc: ""
`,
		`Type: string
Name: "Reserved"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[json:reserved]
Doc: "Shadowing field.\n"
`,
		`Type: string
Name: "A"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: "Field A\n"
`,
		`Type: int
Name: "B"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: "Field B\n"
`,
		`Type: github.com/modulo-srl/mu-config/go2cfg/testdata.EmptySubType
Name: "Test1"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: ""
`,
		`Type: []github.com/modulo-srl/mu-config/go2cfg/testdata.EmptySubType
Name: "Test2"
Layout: 1
Element type: github.com/modulo-srl/mu-config/go2cfg/testdata.EmptySubType
IsEmbedded: false
Tags: map[]
Doc: ""
`,
		`Type: string
Name: "Name"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: "Name describes the protocol name.\nMultiple line documentation test.\nProtocol name.\n"
`,
		`Type: int
Name: "Major"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: "Major version.\n"
`,
		`Type: int
Name: "Minor"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: "Minor version.\n"
`,
		`Type: string
Name: "IP"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: "Remote IP address.\n"
`,
		`Type: int
Name: "Port"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: "Remote port.\n"
`,
		`Type: github.com/modulo-srl/mu-config/go2cfg/testdata.Protocol
Name: "Default"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[json:default_proto]
Doc: "Default protocol.\n"
`,
		`Type: []github.com/modulo-srl/mu-config/go2cfg/testdata.Protocol
Name: "Optionals"
Layout: 1
Element type: github.com/modulo-srl/mu-config/go2cfg/testdata.Protocol
IsEmbedded: false
Tags: map[json:optional_protos]
Doc: "Optional supported protocols.\n"
`,
		`Type: string
Name: "Name"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: "Name of the user documentation block.\nUser name comment.\n"
`,
		`Type: string
Name: "Surname"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: "User surname comment.\n"
`,
		`Type: int
Name: "Age"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[json:age]
Doc: "Age documentation block.\nUser age.\n"
`,
		`Type: int
Name: "StarsCount"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[json:stars_count]
Doc: "Number of stars achieved.\n"
`,
		`Type: []string
Name: "Addresses"
Layout: 1
Element type: string
IsEmbedded: false
Tags: map[]
Doc: "Addresses comment.\n"
`,
		`Type: map[string]string
Name: "Tags"
Layout: 2
Element type: string
IsEmbedded: false
Tags: map[]
Doc: "User tags.\n"
`,
		`Type: github.com/modulo-srl/mu-config/go2cfg/testdata.ConstType
Name: "Type"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: "Type documentation block.\nType of constant.\n"
`,
	}

	whitespacesReplacer := strings.NewReplacer(" ", "◦", "\t", "———➞", "\n", "⏎\n")
	for i, fieldInfo := range info {
		s := fieldInfo.String()
		if s != want[i] {
			t.Fatalf("String return mismatch: got:\n%s\nwant:\n%s\n",
				whitespacesReplacer.Replace(s),
				whitespacesReplacer.Replace(want[i]))
		}
	}
}

func TestFieldInfo(t *testing.T) {
	dirs := []string{"../testdata"}
	testFieldInfo(t, dirs, []*FieldInfoMatch{
		// testdata/embedding.go
		{
			Type: "int", Name: "Identifier",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "id"},
			Doc:  "Identifier documentation block.\n",
		},
		{
			Type: "bool", Name: "Enabled",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Enabled comment line.\n",
		},
		{
			Type: "uint32", Name: "Reserved",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "reserved"},
			Doc:  "",
		},
		{
			Type: "github.com/modulo-srl/mu-config/go2cfg/testdata.Embedded", Name: "",
			Layout: LayoutSingle, IsEmbedded: true,
			Tags: nil,
			Doc:  "Embedded documentation block.\n",
		},
		{
			Type: "float32", Name: "Position",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "position"},
			Doc:  "Position comment line.\n",
		},
		{
			Type: "float32", Name: "Velocity",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "velocity"},
			Doc:  "Velocity documentation block.\n",
		},
		{
			Type: "float32", Name: "Acceleration",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "accel"},
			Doc:  "",
		},
		{
			Type: "string", Name: "Reserved",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "reserved"},
			Doc:  "Shadowing field.\n",
		},
		// testdata/empty.go
		// testdata/empty_defs.go
		{
			Type: "string", Name: "A",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Field A\n",
		},
		{
			Type: "int", Name: "B",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Field B\n",
		},
		{
			Type: "github.com/modulo-srl/mu-config/go2cfg/testdata.EmptySubType", Name: "Test1",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "",
		},
		{
			Type: "[]github.com/modulo-srl/mu-config/go2cfg/testdata.EmptySubType", Name: "Test2",
			Layout: LayoutArray, IsEmbedded: false,
			Tags: nil,
			Doc:  "",
		},
		// testdata/nesting.go
		{
			Type: "string", Name: "Name",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Name describes the protocol name.\nMultiple line documentation test.\nProtocol name.\n",
		},
		{
			Type: "int", Name: "Major",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Major version.\n",
		},
		{
			Type: "int", Name: "Minor",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Minor version.\n",
		},
		{
			Type: "string", Name: "IP",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Remote IP address.\n",
		},
		{
			Type: "int", Name: "Port",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Remote port.\n",
		},
		{
			Type: "github.com/modulo-srl/mu-config/go2cfg/testdata.Protocol", Name: "Default",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "default_proto"},
			Doc:  "Default protocol.\n",
		},
		{
			Type: "[]github.com/modulo-srl/mu-config/go2cfg/testdata.Protocol", Name: "Optionals",
			Layout: LayoutArray, IsEmbedded: false,
			Tags: map[string]string{"json": "optional_protos"},
			Doc:  "Optional supported protocols.\n",
		},
		// testdata/simple.go
		{
			Type: "string", Name: "Name",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Name of the user documentation block.\nUser name comment.\n",
		},
		{
			Type: "string", Name: "Surname",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "User surname comment.\n",
		},
		{
			Type: "int", Name: "Age",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "age"},
			Doc:  "Age documentation block.\nUser age.\n",
		},
		{
			Type: "int", Name: "StarsCount",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "stars_count"},
			Doc:  "Number of stars achieved.\n",
		},
		{
			Type: "[]string", Name: "Addresses",
			Layout: LayoutArray, IsEmbedded: false,
			Tags: nil,
			Doc:  "Addresses comment.\n",
		},
		{
			Type: "map[string]string", Name: "Tags",
			Layout: LayoutMap, IsEmbedded: false,
			Tags: nil,
			Doc:  "User tags.\n",
		},
		{
			Type: "github.com/modulo-srl/mu-config/go2cfg/testdata.ConstType", Name: "Type",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Type documentation block.\nType of constant.\n",
		},
	})
}

func TestFieldInfoMultiPackage(t *testing.T) {
	dirs := []string{
		"../testdata/multipkg/network",
		"../testdata/multipkg/stats",
		"../testdata/multipkg",
	}
	testFieldInfo(t, dirs, []*FieldInfoMatch{
		// testdata/multipkg/status.go
		{
			Type: "bool", Name: "Connected",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Connected flag comment.\n",
		},
		{
			Type:   "github.com/modulo-srl/mu-config/go2cfg/testdata/multipkg/network.ConnState",
			Name:   "State",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Connection state comment.\n",
		},
		// testdata/multipkg/info.go
		{
			Type: "int", Name: "PacketLoss",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "packet_loss"},
			Doc:  "PacketLoss documentation block.\nPacket loss comment.\n",
		},
		{
			Type: "int", Name: "RoundTripTime",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "round_trip_time"},
			Doc:  "Round-trip time in milliseconds.\n",
		},
		// testdata/multipkg/multi_package.go
		{
			Type:   "github.com/modulo-srl/mu-config/go2cfg/testdata/multipkg/network.Status",
			Name:   "NetStatus",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Network status.\n",
		},
		{
			Type:   "github.com/modulo-srl/mu-config/go2cfg/testdata/multipkg/stats.Info",
			Name:   "",
			Layout: LayoutSingle, IsEmbedded: true,
			Tags: nil,
			Doc:  "Statistics info.\n",
		},
	})
}

func testFieldInfo(t *testing.T, patterns []string, want []*FieldInfoMatch) {
	fields := testutils.GetFieldsInfo(t, patterns)

	if len(fields) != len(want) {
		t.Fatalf("Parsed %d fields, want %d.", len(fields), len(want))
	}

	for i, field := range fields {
		if field.Type.String() != want[i].Type || field.Name != want[i].Name ||
			field.Layout != want[i].Layout || field.IsEmbedded != want[i].IsEmbedded ||
			!reflect.DeepEqual(field.Tags, want[i].Tags) || field.Doc != want[i].Doc {
			t.Fatalf("Parsed field mismatch:\n%s\n\nwant:\n%s\n", field, want[i])
		}
	}
}
