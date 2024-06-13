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

const version = "1.0.8"

func main() {
	flag.Usage = usage

	typeName := flag.String("type", "", "struct type name for which generate config; mandatory")
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

	var code string
	var err error

	if *output == "" {
		code, err = generateToml(dir, *typeName, docMode)
		if err != nil {
			log.Fatal(err)
		}

		print(code)
		os.Exit(0)
	}

	ext := filepath.Ext(*output)
	switch ext {
	case ".json":
		fallthrough
	case ".jsonc":
		err = generateJsoncFile(dir, *typeName, *output, docMode)
		if err != nil {
			log.Fatal(err)
		}

	case ".toml":
		err = generateTomlFile(dir, *typeName, *output, docMode)
		if err != nil {
			log.Fatal(err)
		}

	case ".yaml":
		err = generateYamlFile(dir, *typeName, *output, docMode)
		if err != nil {
			log.Fatal(err)
		}

	default:
		err = generateJsoncFile(dir, *typeName, *output+".jsonc", docMode)
		if err != nil {
			log.Fatal(err)
		}

		err = generateTomlFile(dir, *typeName, *output+".toml", docMode)
		if err != nil {
			log.Fatal(err)
		}

		err = generateYamlFile(dir, *typeName, *output+".yaml", docMode)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func usage() {
	println("go2cfg v" + version)

	println("Usage:")
	println("  go2cfg -type <type-name> [-doc-types bits] [-out jsonc-filename] [package-dir]\n")
	println()

	flag.PrintDefaults()

	println("\npackage-dir: directory that contains the go file where specified type is")
	println("defined; when omitted, current working directory will be used\n")
}

func generateJsoncFile(dir, typeName, filename string, docMode renderers.DocTypesMode) error {
	renderer := renderers.NewJsonc(docMode)
	output, err := generator.Generate(dir, typeName, renderer)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, []byte(output), 0666)
	if err != nil {
		return fmt.Errorf("generating JSON file: %s", err)
	}

	return nil
}

func generateTomlFile(dir, typeName, filename string, docMode renderers.DocTypesMode) error {
	output, err := generateToml(dir, typeName, docMode)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, []byte(output), 0666)
	if err != nil {
		return fmt.Errorf("generating TOML file: %s", err)
	}

	return nil
}

func generateToml(dir, typeName string, docMode renderers.DocTypesMode) (string, error) {
	renderer := renderers.NewToml(docMode, true)
	output, err := generator.Generate(dir, typeName, renderer)
	if err != nil {
		return "", fmt.Errorf("generating TOML: %s", err)
	}

	return output, nil
}

func generateYamlFile(dir, typeName, filename string, docMode renderers.DocTypesMode) error {
	renderer := renderers.NewYaml(docMode, 2)
	output, err := generator.Generate(dir, typeName, renderer)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, []byte(output), 0666)
	if err != nil {
		return fmt.Errorf("generating YAML file: %s", err)
	}

	return nil
}
