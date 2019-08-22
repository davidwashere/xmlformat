package goxmlformat

import (
	"fmt"
	"strings"
)

func identStr(indent int) string {
	return strings.Repeat("   ", indent)
}

// FormatXML - Returns xmlStr formatted for Pretty Printing
func FormatXML(xmlStr string) string {
	var final strings.Builder
	var buffer strings.Builder
	var rolling string

	prevFinished := false
	hitNiner := false
	inCDATA := false

	var indent int

	for _, c := range xmlStr {
		buffer.WriteRune(c)

		if hitNiner {
			rolling = rolling[1:] + string(c)
		} else {
			rolling += string(c)
			if len(rolling) == 9 {
				hitNiner = true
			}
		}

		if inCDATA {
			if strings.HasSuffix(rolling, "]]>") {
				inCDATA = false
			}
			continue
		}

		if rolling == "<![CDATA[" {
			inCDATA = true
			continue
		}

		if c == '>' {
			bufStr := buffer.String()
			first := strings.LastIndex(bufStr, "<") + 1
			last := len(bufStr) - 2

			if bufStr[first] == '/' {
				// handles 'end tags' </end>
				indent--
				if prevFinished {
					bufStr = strings.TrimSpace(bufStr)
					fmt.Fprintf(&final, "%s%s\n", identStr(indent), bufStr)
				} else {
					fmt.Fprintf(&final, "%s\n", bufStr)
				}
				prevFinished = true

			} else if bufStr[first] == '?' || bufStr[last] == '/' || bufStr[first] == '!' {
				// handles header <?xml ... ?>, self closing tags <br />, and comments <!-- blah -->
				bufStr = strings.TrimSpace(bufStr)

				if prevFinished {
					fmt.Fprintf(&final, "%s%s\n", identStr(indent), bufStr)
				} else {
					fmt.Fprintf(&final, "\n%s%s\n", identStr(indent), bufStr)
				}
				prevFinished = true

			} else {
				// handles start tags <start>
				bufStr = strings.TrimSpace(bufStr)
				if prevFinished {
					fmt.Fprintf(&final, "%s%s", identStr(indent), bufStr)
				} else {
					fmt.Fprintf(&final, "\n%s%s", identStr(indent), bufStr)
				}
				prevFinished = false
				indent++
			}

			buffer.Reset()
		}
	}

	return strings.TrimSpace(final.String())
}
