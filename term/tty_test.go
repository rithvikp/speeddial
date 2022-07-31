package term

import (
	"testing"
)

func TestNumLinesInWidth(t *testing.T) {
	tests := []struct {
		msg   string
		input string
		width int
		lines int
	}{
		{
			msg:   "Empty string",
			input: "",
			width: 100,
			lines: 0,
		},
		{
			msg:   "ASCII single-line string",
			input: "first line",
			width: 100,
			lines: 1,
		},
		{
			msg:   "ASCII multi-line string",
			input: "first line\nsecond line\nthird line",
			width: 100,
			lines: 3,
		},
		{
			msg:   "Multi-line string with escape sequences",
			input: "first line\nsecond line                                    \x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[90m\x1b[90m | \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39mA separate column on the same line                \x1b[0m\x1b[0m\x1b[K\r\n\x1b[39m\x1b[39mthird line",
			width: 100,
			lines: 3,
		},
		{
			msg:   "Multi-line string with escape sequences and wrapping",
			input: "first line\nsecond line                                    \x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[90m\x1b[90m | \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39mA separate column on the same line                \x1b[0m\x1b[0m\x1b[K\r\n\x1b[39m\x1b[39mthird line",
			width: 50,
			lines: 4,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.msg, func(t *testing.T) {
			got := numLinesInWidth(tt.input, tt.width)

			if got != tt.lines {
				t.Errorf("Returned line count was incorrect: got %d, want %d", got, tt.lines)
			}
		})
	}
}
