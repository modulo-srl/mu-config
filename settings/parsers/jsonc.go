package parsers

import (
	"bytes"
	"encoding/json"
	"os"
	"regexp"
)

func LoadJsoncFile(filename string, data interface{}) error {
	bb, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return LoadJsonc(bb, data)
}

func LoadJsonc(bb []byte, data interface{}) error {
	// Converte da Jsonc a Json.
	bb = translate(bb)

	r := bytes.NewReader(bb)

	d := json.NewDecoder(r)
	d.DisallowUnknownFields()

	err := d.Decode(data)
	if err != nil {
		return err
	}

	return nil
}

// translate converte da jsonc a json.
func translate(s []byte) []byte {
	commentsRE := regexp.MustCompile(`/\*(?:.|\n)*?\*/|//[^\r\n]*`)
	stringRE := regexp.MustCompile(`"(?:\\(?:u[0-9a-fA-F]{4}|.)|[^\\"])*"`)

	commentsRanges := commentsRE.FindAllIndex(s, -1)
	// Se non ci sono commenti ritorna la slice originale.
	if len(commentsRanges) == 0 {
		return s
	}

	// Cerca le stringhe per scartare falsi positivi dai match commentsRanges.
	stringsRanges := stringRE.FindAllIndex(s, -1)

	j := make([]byte, len(s))
	copy(j, s)

	var skip bool
	offset := 0
	firstString := 0
	for _, comment := range commentsRanges {
		skip = false
		for i := firstString; i < len(stringsRanges); i++ {
			str := stringsRanges[i]
			if comment[0] > str[0] && comment[0] < str[1] {
				// Il commento inizia all'interno di una stringa, lo ignora.
				skip = true
				break
			}
			if str[0] > comment[1] {
				// La stringa inizia dopo la fine del commento.
				// Salva l'indice per la prossima ricerca e interrompe il loop.
				firstString = i
				break
			}
		}

		if skip {
			continue
		}

		// Rimuove il commento.
		j = append(j[:comment[0]+offset], j[comment[1]+offset:]...)
		offset -= comment[1] - comment[0]
	}

	return j
}
