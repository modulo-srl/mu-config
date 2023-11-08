package testdata

//go:generate go2cfg -type Embedding -out embedding.jsonc
//go:generate go2cfg -type Embedding -doc-types all -out embedding_all_fields.jsonc
//go:generate go2cfg -type Embedding -doc-types basic -out embedding_basic_fields.jsonc

// Embedded test struct.
type Embedded struct {
	// Identifier documentation block.
	Identifier int  `json:"id"`
	Enabled    bool // Enabled comment line.

	Reserved uint32 `json:"reserved"`
}

// Embedding test struct.
type Embedding struct {
	// Embedded documentation block.
	Embedded

	Position float32 `json:"position"` // Position comment line.
	// Velocity documentation block.
	Velocity     float32 `json:"velocity"`
	Acceleration float32 `json:"accel"`

	Reserved string `json:"reserved"` // Shadowing field.
}

func EmbeddingDefaults() *Embedding {
	return &Embedding{
		Embedded: Embedded{
			Identifier: 1234,
			Enabled:    false,
			Reserved:   0x10,
		},
		Position:     1.0,
		Velocity:     2,
		Acceleration: 0.23,
		Reserved:     "Shadowing",
	}
}
