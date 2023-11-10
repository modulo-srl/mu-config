package renderers

import (
	"github.com/modulo-srl/mu-config/go2cfg/distiller"
	"github.com/modulo-srl/mu-config/go2cfg/ordered"
	"go/types"
)

// DocTypesMode defines rendering modes for field types in JSONC comments.
type DocTypesMode int

const (
	NoFields    DocTypesMode = iota // Don't show type on all fields (default).
	AllFields                       // Show types for all fields.
	BasicFields                     // Show types for basic fields only (int, float, bool, string).
)

// Interface defines a renderer interface.
type Interface interface {
	// RenderStruct renders indented code for specified struct and all nested or embedded ones recursively.
	RenderStruct(info *distiller.StructInfo, defaults interface{}, indent string,
		embedded bool, parentShadowing []string) (string, error)
	// RenderArray renders slice or array fields.
	RenderArray(field *distiller.FieldInfo, value []interface{}, indent string) (string, error)
	// RenderMap renders map fields.
	RenderMap(field *distiller.FieldInfo, value *ordered.Map, indent string) (string, error)
	// RenderElement renders an element value of a slice, array or map.
	RenderElement(itemType types.Type, item interface{}, indent string) (string, error)
}
