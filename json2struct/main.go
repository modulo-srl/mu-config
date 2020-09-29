package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/modulo-srl/mu-settings/json2struct/encoder"
)

const version = "json2struct 1.0"

func usage() {
	println(version)
	println("Usage: json2struct [-p packageName] [-s structName] [-f initFuncName] filename.json")
	flag.PrintDefaults()
}

func main() {
	ver := flag.Bool("version", false, "Show tool version")
	declare := flag.Bool("declare", false, "Make struct declaration")
	values := flag.Bool("values", false, "Fill struct with values")
	pkg := flag.String("pkg", "config", "Package name to generate")
	structName := flag.String("struct", "configData", "Output struct name")
	funcName := flag.String("func", "", `Init function name (default "initStructName")`)
	f64 := flag.Bool("f64", false, "Force 64bit for integer and float")
	w := flag.Bool("warn", false, "Generate warning header")
	outFile := flag.String("out", "", "Save to file instead of stdout")
	flag.Parse()

	var filename string

	if *ver {
		fmt.Println(version)
		return
	}

	if flag.NArg() == 0 || (!*declare && !*values) {
		usage()
		return
	}

	//filename = "../../bin/settings.default.json"
	filename = flag.Arg(0)
	out := ""

	out += "package " + *pkg + "\n\n"

	if *w {
		out += "// WARNING: autogenerated file by " + version + ". Do not edit.\n" +
			"// cmdline: `" + strings.Join(os.Args, " ") + "`\n\n"
	}

	if *declare {
		decoded, err := encoder.JSONToStruct(filename, *structName, *f64)
		if err != nil {
			panic(err)
		}
		out += decoded
	}

	if *values {
		decoded, err := encoder.JSONToValues(filename, *funcName, *structName, *f64)
		if err != nil {
			panic(err)
		}
		out += decoded
	}

	if *outFile != "" {
		err := encoder.WriteFile(*outFile, out)
		if err != nil {
			panic(err)
		}
		return
	}

	fmt.Println(out)
}
