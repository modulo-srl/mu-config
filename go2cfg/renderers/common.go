package renderers

// Functions common to all renderers.

import (
	"fmt"
	"github.com/modulo-srl/mu-config/go2cfg/distiller"
	"go/constant"
	"go/types"
	"log"
	"strings"
)

// renderDoc renders formatted field documentation indenting it with passed indent string,
// starting with marker string.
func renderDoc(f *distiller.FieldInfo, indent string, marker string, renderType bool) string {
	doc := f.Doc

	// Check if the type is used to define typed constants.
	consts := distiller.LookupTypedConsts(f.Type.String())
	if consts != nil {
		// Display allowed values for defined constants below the field documentation.
		doc += "Allowed values:\n"

		constLen := 0
		valueLen := 0
		for _, info := range consts {
			if len(info.Name) > constLen {
				constLen = len(info.Name)
			}
			if len(info.Value) > valueLen {
				valueLen = len(info.Value)
			}
		}

		for _, info := range consts {
			doc += fmt.Sprintf("%-*s = %*v  %s\n", constLen, info.Name, valueLen, info.Value, info.InlineDoc())
		}
	}

	// Indent the documentation.
	commentPrefix := indent + marker + " "
	d := strings.ReplaceAll(doc, "\n", "\n"+commentPrefix)
	if len(d) > 0 {
		d = d[:len(d)-len(commentPrefix)]
	} else {
		d = "\n"
	}

	typeName := ""

	if renderType {
		typeName = f.Type.String()
		if lastSlash := strings.LastIndex(typeName, "/"); lastSlash >= 0 {
			typeName = typeName[lastSlash+1:]
			// The square brackets at the beginning of the typeName are trimmed out, so must be re-added.
			if f.Layout == distiller.LayoutArray {
				typeName = "[]" + typeName
			}
		}

		if d != "\n" {
			typeName += " - "
		}

		d = typeName + d
	}

	if d == "\n" {
		return ""
	}

	return commentPrefix + d
}

// lastIndexOf returns the last slice index of specified value.
func lastIndexOf(slice []string, value string) int {
	if slice != nil {
		for i := len(slice) - 1; i >= 0; i-- {
			if slice[i] == value {
				return i
			}
		}
	}

	return -1
}

// typeZero return the default uninitialized value for specified field.
func typeZero(field *distiller.FieldInfo) interface{} {
	var value interface{}
	if field.Layout == distiller.LayoutArray {
		value = make([]interface{}, 0)
		return value
	} else if field.Layout == distiller.LayoutMap {
		value = make(map[interface{}]interface{})
		return value
	}

	fieldType := types.Default(field.Type)
	switch t := fieldType.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Bool:
			value = *new(bool)
		case types.Int:
			value = *new(int)
		case types.Int8:
			value = *new(int8)
		case types.Int16:
			value = *new(int16)
		case types.Int32:
			value = *new(int32)
		case types.Int64:
			value = *new(int64)
		case types.Uint:
			value = *new(uint)
		case types.Uint8:
			value = *new(uint8)
		case types.Uint16:
			value = *new(uint16)
		case types.Uint32:
			value = *new(uint32)
		case types.Uint64:
			value = *new(uint64)
		case types.Uintptr:
			value = *new(uintptr)
		case types.Float32:
			value = *new(float32)
		case types.Float64:
			value = *new(float64)
		case types.Complex64:
			value = *new(complex64)
		case types.Complex128:
			value = *new(complex128)
		case types.String:
			value = constant.MakeString("")
		default:
			log.Fatalf("Unhandled default value for type %v", fieldType.String())
		}
	}

	return value
}

// fieldsSlice defines a slice of fields.
type fieldsSlice []*distiller.FieldInfo

// indexOf return index of field into the slice, -1 if not found.
func (fs fieldsSlice) indexOf(name string) int {
	for i, field := range fs {
		if field.Name == name {
			return i
		}
	}

	return -1
}
