// Package generator generates code from the information extracted
// from the AST via the distiller package.
package generator

import (
	"fmt"
	"github.com/modulo-srl/mu-config/go2cfg/distiller"
	"github.com/modulo-srl/mu-config/go2cfg/renderers"
	"log"
)

// Generate generates code for given package dir and type name.
func Generate(dir, typeName string, renderer renderers.Interface) (string, error) {
	pkgInfo, err := distiller.NewPackageInfo(dir, typeName)
	if err != nil {
		return "", err
	}

	s := distiller.LookupStruct(pkgInfo.Package.PkgPath + "." + typeName)
	if s == nil {
		return "", fmt.Errorf("cannot find struct %s in package %s", typeName, pkgInfo.Package.Name)
	}

	var code string
	code, err = renderer.RenderStruct(s, s.Defaults, "", false, nil)
	if err != nil {
		log.Fatal(err)
	}

	return code, nil
}
