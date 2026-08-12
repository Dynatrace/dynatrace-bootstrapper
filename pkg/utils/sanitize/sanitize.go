package sanitize

import "strings"

var controlCharReplacer = strings.NewReplacer("\n", "", "\t", "", "\r", "", "\x00", "")

// StripControlChars removes control characters (newline, carriage return, tab, null) that
// could be used to inject forged directives or entries into INI/properties files.
func StripControlChars(s string) string {
	return controlCharReplacer.Replace(s)
}
