package renderers

import (
	"fmt"
	"go/constant"
	"regexp"
)

var unescapeStringRegexp = regexp.MustCompile(`\\(.)`)

// unescapeString unescapes a literal string value representation.
func unescapeString(v any) string {
	s := fmt.Sprintf("%v", v)
	if val, ok := v.(constant.Value); ok {
		s = val.ExactString()
	}

	if len(s) > 0 {
		s = unescapeStringRegexp.ReplaceAllStringFunc(s[1:len(s)-1], func(s string) string {
			switch s[1:] {
			case "r":
				return "\r"
			case "n":
				return "\n"
			case "t":
				return "\t"
			case "f":
				return "\f"
			case "v":
				return "\v"
			}

			return s[1:]
		})
	}

	return s
}
