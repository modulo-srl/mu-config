package distiller

import "testing"

func TestPackageInfo(t *testing.T) {
	info, err := NewPackageInfo("../testdata", "")
	if err != nil {
		t.Fatal(err)
	}

	wantConsts := []string{
		"github.com/modulo-srl/mu-config/go2cfg/testdata.ConstType",
	}

	wantStructs := []string{
		"github.com/modulo-srl/mu-config/go2cfg/testdata.Embedded",
		"github.com/modulo-srl/mu-config/go2cfg/testdata.Embedding",
		"github.com/modulo-srl/mu-config/go2cfg/testdata.Empty",
		"github.com/modulo-srl/mu-config/go2cfg/testdata.EmptyDefs",
		"github.com/modulo-srl/mu-config/go2cfg/testdata.EmptySubType",
		"github.com/modulo-srl/mu-config/go2cfg/testdata.Protocol",
		"github.com/modulo-srl/mu-config/go2cfg/testdata.Nesting",
		"github.com/modulo-srl/mu-config/go2cfg/testdata.Simple",
	}

	if len(info.TypedConsts) != len(wantConsts) {
		t.Fatalf("Parsed %d types used by typed constants, want %d",
			len(info.TypedConsts), len(wantConsts))
	}

	if len(info.Structs) != len(wantStructs) {
		t.Fatalf("Parsed %d structs, want %d", len(info.Structs), len(wantStructs))
	}

	for _, name := range wantConsts {
		if LookupTypedConsts(name) == nil {
			t.Fatalf("Cannot lookup typed constants of type %s", name)
		}
	}

	for _, name := range wantStructs {
		if LookupStruct(name) == nil {
			t.Fatalf("Cannot lookup struct %s", name)
		}
	}

	_, err = NewPackageInfo("../invalid-path", "")
	if err == nil {
		t.Fatal("Lookup for invalid package path, error expected, got nil.")
	}

	_, err = NewPackageInfo("../testdata/consts.go", "")
	if err == nil {
		t.Fatal("Non directory package path, error expected, got nil.")
	}
}

func TestPackageInfoMultiPackage(t *testing.T) {
	_, err := NewPackageInfo("../testdata/multipkg", "")
	if err != nil {
		t.Fatal(err)
	}

	wantConsts := []string{
		"github.com/modulo-srl/mu-config/go2cfg/testdata/multipkg/network.ConnState",
	}

	wantStructs := []string{
		"github.com/modulo-srl/mu-config/go2cfg/testdata/multipkg/network.Status",
		"github.com/modulo-srl/mu-config/go2cfg/testdata/multipkg/stats.Info",
		"github.com/modulo-srl/mu-config/go2cfg/testdata/multipkg.MultiPackage",
	}

	for _, name := range wantConsts {
		if LookupTypedConsts(name) == nil {
			t.Fatalf("Cannot lookup typed constants of type %s", name)
		}
	}

	if LookupTypedConsts("invalid-name") != nil {
		t.Fatalf("Lookup of invalid typed constant not nil")
	}

	for _, name := range wantStructs {
		if LookupStruct(name) == nil {
			t.Fatalf("Cannot lookup struct %s", name)
		}
	}

	if LookupStruct("invalid-name") != nil {
		t.Fatalf("Lookup of invalid struct not nil")
	}

	_, err = NewPackageInfo("../testdata/multipkg", "MultiPackage")
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewPackageInfo("../testdata/multipkg", "invalid-name")
	if err == nil {
		t.Fatalf("Lookup of invalid struct, error expected, got nil")
	}
}
