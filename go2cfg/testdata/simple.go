package testdata

//go:generate go2cfg -type Simple -out simple.jsonc
//go:generate go2cfg -type Simple -doc-types all -out simple_all_fields.jsonc
//go:generate go2cfg -type Simple -doc-types basic -out simple_basic_fields.jsonc

// Simple defines a simple user.
type Simple struct {
	// Name of the user documentation block.
	Name    string // User name comment.
	Surname string // User surname comment.

	// Age documentation block.
	Age        int `json:"age"`         // User age.
	StarsCount int `json:"stars_count"` // Number of stars achieved.

	Addresses []string // Addresses comment.

	Tags map[string]string // User tags.

	// Type documentation block.
	Type ConstType // Type of constant.

	// X, Y documentation block.
	X, Y float64 // Coordinates.
}

func SimpleDefaults() *Simple {
	return &Simple{
		Name:       "John",
		Age:        30,
		StarsCount: 5,
		Addresses: []string{
			"Address 1",
			"Address 2",
			"Address 3",
		},
		Tags: map[string]string{
			"Key1": "Value1",
			"Key2": "Value2",
			"Key3": "Value3",
		},
		X: 1.0,
		Y: 2.0,
	}
}
