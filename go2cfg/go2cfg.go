package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/modulo-srl/mu-config/go2cfg/generator"
	"github.com/modulo-srl/mu-config/go2cfg/renderers"
)

const version = "1.0"

func main() {
	flag.Usage = usage

	typeName := flag.String("type", "", "struct type name for which generate JSONC; mandatory")
	output := flag.String("out", "", "output filepath; The extension in the filepath\n"+
		"establishes the type of format to be generated, otherwise without extension\n"+
		"a file for each format will be exported (yaml, toml, jsonc).\n"+
		"When omitted outputs to stdout in toml format")
	docTypeMode := flag.String("doc-types", "",
		"Struct fields types for which render the type in comments:\n"+
			"  all    Display type in all fields;\n"+
			"  basic  Display type in fields of basic type (int, float, bool, string).")

	flag.Parse()

	if *typeName == "" {
		println("Flag -type is mandatory.\n")
		flag.Usage()
		os.Exit(1)
	}

	docMode := renderers.NoFields
	if *docTypeMode != "" {
		switch *docTypeMode {
		case "all":
			docMode = renderers.AllFields
		case "basic":
			docMode = renderers.BasicFields
		default:
			fmt.Printf("Invalid value \"%s\" for -doc-types flag.\n\n", *docTypeMode)
			flag.Usage()
			os.Exit(1)
		}
	}

	dirs := flag.Args()

	dir := "."
	switch len(dirs) {
	case 0:
		println("No directory specified, using current working dir.")

	case 1:
		dir = dirs[0]

	default:
		println("Only one directory can be specified.\n")
		flag.Usage()
		os.Exit(1)
	}

	var renderer renderers.Interface
	ext := filepath.Ext(*output)
	switch ext {
	case ".jsonc":
		renderer = renderers.NewJsonc(docMode)

	case ".toml":
		renderer = renderers.NewToml(docMode, false)

	default:
		log.Fatalf("unsupported output format: %s", ext)
	}

	code, err := generator.Generate(dir, *typeName, renderer)
	if err != nil {
		log.Fatal(err)
	}

	if *output != "" {
		err = os.WriteFile(*output, []byte(code), 0666)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		print(code)
	}
}

func usage() {
	println("go2cfg v" + version)

	println("Usage:")
	println("  go2cfg -type <type-name> [-doc-types bits] [-out jsonc-filename] [package-dir]\n")

	flag.PrintDefaults()

	println("\npackage-dir: directory that contains the go file where specified type is")
	println("defined; when omitted, current working directory will be used\n")

	println("Allowed constants for -doc-types flag:")
	println("  all    Display type in all fields;")
	println("  basic  Display type in fields of basic type (int, float, bool, string).")
}
