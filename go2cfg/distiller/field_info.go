package distiller

import (
	"fmt"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/packages"
	"log"
	"regexp"
	"strings"
)

type FieldLayout int

const (
	LayoutSingle FieldLayout = iota // The field is a single element.
	LayoutArray                     // The field is an array or slice of elements.
	LayoutMap                       // The field is a map of elements.
)

// FieldInfo holds information about structure field.
type FieldInfo struct {
	Type       types.Type        // Field type, used to compute fully qualified type string.
	Name       string            // Field name.
	Layout     FieldLayout       // Field layout.
	EltType    types.Type        // Field element type, when the field is a slice or map.
	IsEmbedded bool              // True if field is an embedded struct and Name is an empty string.
	Tags       map[string]string // Tags applied to that field as map of name-value key-pairs.
	Doc        string            // Documentation content if present.
}

// tagRegexp defines a regex to extract tags names and values.
var tagRegexp = regexp.MustCompile(`(\w+):"((?:[^"\\]|\\.)*)"`)

// NewFieldInfo creates new field information object from given abstract syntax tree field and package.
// Terminates the process with a fatal error if multiple names are specified for the same field.
func NewFieldInfo(field *ast.Field, pkg *packages.Package) *FieldInfo {
	f := &FieldInfo{Layout: LayoutSingle, EltType: nil}
	if len(field.Names) == 1 {
		f.Name = field.Names[0].Name
	} else if field.Names == nil {
		// Embedded field.
		f.IsEmbedded = true
	} else {
		log.Fatalf("Unsupported multiple names.")
	}

	f.Type = pkg.TypesInfo.Types[field.Type].Type
	switch fieldType := field.Type.(type) {
	case *ast.ArrayType:
		// In case of array get the type of single element.
		f.EltType = pkg.TypesInfo.Types[fieldType.Elt].Type
		f.Layout = LayoutArray

	case *ast.MapType:
		// In case of map get the type of value.
		f.EltType = pkg.TypesInfo.Types[fieldType.Value].Type
		f.Layout = LayoutMap
	}

	// Parse defined tags populating FieldInfo.Tags map.
	if field.Tag != nil {
		f.Tags = make(map[string]string)
		tags := tagRegexp.FindAllStringSubmatch(strings.Trim(field.Tag.Value, "` "), -1)
		for _, tag := range tags {
			tagValue := ""
			if len(tag) == 3 {
				tagValue = tag[2]
			}
			f.Tags[tag[1]] = tagValue
		}
	}

	// Merge documentation and comment.
	f.Doc = field.Doc.Text() + field.Comment.Text()
	return f
}

func (f *FieldInfo) String() string {
	return fmt.Sprintf("Type: %s\nName: \"%s\"\nLayout: %v\nElement type: %v\nIsEmbedded: %v\nTags: %+v\nDoc: \"%v\"\n",
		f.Type.String(), f.Name, f.Layout, f.EltType,
		f.IsEmbedded, f.Tags, strings.ReplaceAll(f.Doc, "\n", "\\n"))
}
