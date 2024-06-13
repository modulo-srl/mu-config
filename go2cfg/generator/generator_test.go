package generator

import (
	"github.com/modulo-srl/mu-config/go2cfg/renderers"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerator(t *testing.T) {
	var tests = []struct {
		pkgDir   string
		typeName string
		filename string
		mode     renderers.DocTypesMode
	}{
		// jsonc
		{"../testdata", "Embedding", "../testdata/embedding.jsonc", renderers.NoFields},
		{"../testdata", "Empty", "../testdata/empty.jsonc", renderers.NoFields},
		{"../testdata", "Nesting", "../testdata/nesting.jsonc", renderers.NoFields},
		{"../testdata", "Simple", "../testdata/simple.jsonc", renderers.NoFields},
		{"../testdata/multipkg", "MultiPackage", "../testdata/multipkg/multi_package.jsonc", renderers.NoFields},

		{"../testdata", "Embedding", "../testdata/embedding_basic_fields.jsonc", renderers.BasicFields},
		{"../testdata", "Nesting", "../testdata/nesting_basic_fields.jsonc", renderers.BasicFields},
		{"../testdata", "Simple", "../testdata/simple_basic_fields.jsonc", renderers.BasicFields},
		{"../testdata/multipkg", "MultiPackage", "../testdata/multipkg/multi_package_basic_fields.jsonc", renderers.BasicFields},

		{"../testdata", "Embedding", "../testdata/embedding_all_fields.jsonc", renderers.AllFields},
		{"../testdata", "Nesting", "../testdata/nesting_all_fields.jsonc", renderers.AllFields},
		{"../testdata", "Simple", "../testdata/simple_all_fields.jsonc", renderers.AllFields},
		{"../testdata/multipkg", "MultiPackage", "../testdata/multipkg/multi_package_all_fields.jsonc", renderers.AllFields},

		// toml
		{"../testdata", "Embedding", "../testdata/embedding.toml", renderers.NoFields},
		{"../testdata", "Empty", "../testdata/empty.toml", renderers.NoFields},
		{"../testdata", "Nesting", "../testdata/nesting.toml", renderers.NoFields},
		{"../testdata", "Simple", "../testdata/simple.toml", renderers.NoFields},
		{"../testdata/multipkg", "MultiPackage", "../testdata/multipkg/multi_package.toml", renderers.NoFields},

		{"../testdata", "Embedding", "../testdata/embedding_basic_fields.toml", renderers.BasicFields},
		{"../testdata", "Nesting", "../testdata/nesting_basic_fields.toml", renderers.BasicFields},
		{"../testdata", "Simple", "../testdata/simple_basic_fields.toml", renderers.BasicFields},
		{"../testdata/multipkg", "MultiPackage", "../testdata/multipkg/multi_package_basic_fields.toml", renderers.BasicFields},

		{"../testdata", "Embedding", "../testdata/embedding_all_fields.toml", renderers.AllFields},
		{"../testdata", "Nesting", "../testdata/nesting_all_fields.toml", renderers.AllFields},
		{"../testdata", "Simple", "../testdata/simple_all_fields.toml", renderers.AllFields},
		{"../testdata/multipkg", "MultiPackage", "../testdata/multipkg/multi_package_all_fields.toml", renderers.AllFields},

		// yaml
		{"../testdata", "Embedding", "../testdata/embedding.yaml", renderers.NoFields},
		{"../testdata", "Empty", "../testdata/empty.yaml", renderers.NoFields},
		{"../testdata", "Nesting", "../testdata/nesting.yaml", renderers.NoFields},
		{"../testdata", "Simple", "../testdata/simple.yaml", renderers.NoFields},
		{"../testdata/multipkg", "MultiPackage", "../testdata/multipkg/multi_package.yaml", renderers.NoFields},

		{"../testdata", "Embedding", "../testdata/embedding_basic_fields.yaml", renderers.BasicFields},
		{"../testdata", "Nesting", "../testdata/nesting_basic_fields.yaml", renderers.BasicFields},
		{"../testdata", "Simple", "../testdata/simple_basic_fields.yaml", renderers.BasicFields},
		{"../testdata/multipkg", "MultiPackage", "../testdata/multipkg/multi_package_basic_fields.yaml", renderers.BasicFields},

		{"../testdata", "Embedding", "../testdata/embedding_all_fields.yaml", renderers.AllFields},
		{"../testdata", "Nesting", "../testdata/nesting_all_fields.yaml", renderers.AllFields},
		{"../testdata", "Simple", "../testdata/simple_all_fields.yaml", renderers.AllFields},
		{"../testdata/multipkg", "MultiPackage", "../testdata/multipkg/multi_package_all_fields.yaml", renderers.AllFields},
	}

	whitespacesReplacer := strings.NewReplacer(" ", "◦", "\t", "———➞", "\n", "⏎\n")
	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			var renderer renderers.Interface
			switch filepath.Ext(test.filename) {
			case ".jsonc":
				renderer = renderers.NewJsonc(test.mode)

			case ".toml":
				renderer = renderers.NewToml(test.mode, true)

			case ".yaml":
				renderer = renderers.NewYaml(test.mode, 2)

			default:
				t.Fatalf("unsupported file format: %s", test.filename)
			}

			code, err := Generate(test.pkgDir, test.typeName, renderer)
			if err != nil {
				t.Fatal(err)
			}

			content, err := os.ReadFile(test.filename)
			if err != nil {
				t.Fatal(err)
			}

			want := string(content)

			if code != want {
				t.Fatalf("Generated code mismatch for %s struct:\n%s\n\nwant %s:\n%s",
					test.typeName,
					whitespacesReplacer.Replace(code),
					test.filename,
					whitespacesReplacer.Replace(want))
			}
		})
	}

	rr := []renderers.Interface{
		renderers.NewJsonc(renderers.AllFields),
		renderers.NewToml(renderers.AllFields, true),
		renderers.NewYaml(renderers.AllFields, 2),
	}

	for _, r := range rr {
		_, err := Generate("../testdata/invalid-path", "", r)
		if err == nil {
			t.Fatalf("Generating for invalid path: expected error, got nil.")
		}

		_, err = Generate("../testdata", "invalid-struct", r)
		if err == nil {
			t.Fatalf("Generating for invalid struct: expected error, got nil.")
		}
	}
}
