package sanitize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStripControlChars(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{name: "empty string => unchanged", in: "", out: ""},
		{name: "clean string => unchanged", in: "foo=bar", out: "foo=bar"},
		{name: "regular whitespace => kept", in: "foo  bar", out: "foo  bar"},
		{name: "newline => stripped", in: "foo\nbar", out: "foobar"},
		{name: "tab => stripped", in: "foo\tbar", out: "foobar"},
		{name: "carriage return => stripped", in: "foo\rbar", out: "foobar"},
		{name: "null byte => stripped", in: "foo\x00bar", out: "foobar"},
		{name: "multiple control chars at once => all stripped", in: "\nfoo\t\rbar\x00", out: "foobar"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.out, StripControlChars(tt.in))
		})
	}
}
