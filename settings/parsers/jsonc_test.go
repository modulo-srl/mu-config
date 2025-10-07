package parsers

import (
	"strings"
	"testing"
)

func TestTranslate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		// Gestione base e terminatori di riga (LF/CRLF)
		// Nessun commento.
		{"NoComments", `{"key": "value"}`, `{"key": "value"}`, false},
		// Commento linea singola LF.
		{"SingleLineLF", `{"k": "v"} // c` + "\n" + `{"k2": "v2"}`, `{"k": "v"} ` + "\n" + `{"k2": "v2"}`, false},
		// Commento linea singola CRLF.
		{"SingleLineCRLF", `{"k": "v"} // c` + "\r\n" + `{"k2": "v2"}`, `{"k": "v"} ` + "\r\n" + `{"k2": "v2"}`, false},
		// Commento multi-linea.
		{"MultiLineSimple", `/*c1*/{"k": "v"}`, `{"k": "v"}`, false},
		// Commento multi-linea LF.
		{"MultiLineWithLF", "{\n  /* commento\n   su più linee\n  */\n  \"k\": \"v\"\n}", "{\n  \n  \"k\": \"v\"\n}", false},
		// Commento multi-linea CRLF.
		{"MultiLineWithCRLF", "{\r\n  /* commento\r\n   su più linee\r\n  */\r\n  \"k\": \"v\"\r\n}", "{\r\n  \r\n  \"k\": \"v\"\r\n}", false},
		// Commento linea Singola con CRLF e spazi.
		{"CRLFWithSpaces", "  // c \r\n{\"k\": \"v\"}", "  \r\n{\"k\": \"v\"}", false},

		// Gestione stringhe e sequenze di escape
		// Sequenza di commento in stringa.
		{"InStringSlashes", `{"url": "https://example.com//path"}`, `{"url": "https://example.com//path"}`, false},
		// Inizio blocco in stringa.
		{"InStringAsterisk", `{"exp": "x /* y"}`, `{"exp": "x /* y"}`, false},
		// Carattere di escape in stringa.
		{"EscapeCharInString", `{"key": "val\\nue"} // c`, `{"key": "val\\nue"} `, false},
		// Sequenza unicode in stringa.
		{"UnicodeSimple", `{"k": "\u2713"} // tick`, `{"k": "\u2713"} `, false},
		// Backslash in stringa.
		{"BackslashInString", `{"k": "path\\to\\file"} // c`, `{"k": "path\\to\\file"} `, false},
		// Stringa con backslash che fa l'escape dei doppi apici.
		{"EscapedQuote", `{"k": "val\\\"ue"} // c`, `{"k": "val\\\"ue"} `, false},
		// Commento linea singola contenente doppi apici.
		{"SingleLineWithQuotes", `{"k": "v"} // c "commented"` + "\n" + `{"k2": "v2"}`, `{"k": "v"} ` + "\n" + `{"k2": "v2"}`, false},
		// Commento multi-linea contenente doppi apici.
		{"MultiLineWithQuotes", "{\n  /* commento\n   su \"più linee\"\n  */\n  \"k\": \"v\"\n}", "{\n  \n  \"k\": \"v\"\n}", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := translate([]byte(test.input))

			if string(got) != test.want {
				t.Fatalf("got\n%s\nwant\n%s\n",
					strings.ReplaceAll(string(got), "\n", `\n`),
					strings.ReplaceAll(test.want, "\n", `\n`))
			}
		})
	}
}
