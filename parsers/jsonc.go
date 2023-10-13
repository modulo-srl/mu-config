package parsers

import (
	"bytes"
	"encoding/json"
	"os"
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

// Original work: https://github.com/muhammadmuzzammil1998/jsonc/blob/master/translator.go
// MIT License
// Copyright (c) 2019 Muhammad Muzzammil

const (
	ESCAPE   = 92
	QUOTE    = 34
	SPACE    = 32
	TAB      = 9
	NEWLINE  = 10
	ASTERISK = 42
	SLASH    = 47
)

func translate(s []byte) []byte {
	var (
		i       int
		quote   bool
		escaped bool
	)
	j := make([]byte, len(s))
	comment := &commentData{}
	for _, ch := range s {
		if ch == ESCAPE || escaped {
			j[i] = ch
			i++
			escaped = !escaped
			continue
		}
		if ch == QUOTE {
			quote = !quote
		}
		if (ch == SPACE || ch == TAB) && !quote {
			continue
		}
		if ch == NEWLINE {
			if comment.isSingleLined {
				comment.stop()
			}
			continue
		}
		if quote && !comment.startted {
			j[i] = ch
			i++
			continue
		}
		if comment.startted {
			if ch == ASTERISK {
				comment.canEnd = true
				continue
			}
			if comment.canEnd && ch == SLASH {
				comment.stop()
				continue
			}
			comment.canEnd = false
			continue
		}
		if comment.canStart && (ch == ASTERISK || ch == SLASH) {
			comment.start(ch)
			continue
		}
		if ch == SLASH {
			comment.canStart = true
			continue
		}
		j[i] = ch
		i++
	}
	return j[:i]
}

type commentData struct {
	canStart      bool
	canEnd        bool
	startted      bool
	isSingleLined bool
	//endLine       int
}

func (c *commentData) stop() {
	c.startted = false
	c.canStart = false
}

func (c *commentData) start(ch byte) {
	c.startted = true
	c.isSingleLined = ch == SLASH
}
