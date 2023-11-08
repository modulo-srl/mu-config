package testdata

//go:generate go2cfg -type EmptyDefs -out empty_defs.jsonc
//go:generate go2cfg -type EmptyDefs -doc-types all -out empty_all_fields.jsonc
//go:generate go2cfg -type EmptyDefs -doc-types basic -out empty_basic_fields.jsonc

// EmptySubType define a struct with non-initialized fields.
type EmptySubType struct {
	A string // Field A
	B int    // Field B
}

// EmptyDefs define a struct with non-initialized fields.
type EmptyDefs struct {
	Test1 EmptySubType
	Test2 []EmptySubType
}

func EmptyDefsDefaults() *EmptyDefs {
	return &EmptyDefs{}
}
