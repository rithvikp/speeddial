package termui

import (
	"testing"
)

func TestNumNewLinesInWidth(t *testing.T) {
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
			msg:   "Single-line string without trailing new line",
			input: "first line",
			width: 100,
			lines: 0,
		},
		{
			msg:   "Single-line string  withtrailing new line",
			input: "first line",
			width: 100,
			lines: 0,
		},
		{
			msg:   "Multi-line string",
			input: "first line\nsecond line\nthird line",
			width: 100,
			lines: 2,
		},
		{
			msg:   "Multi-line string with color escape sequences",
			input: "first line\nsecond line                                    \x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[90m\x1b[90m | \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39mA separate column on the same line                \x1b[0m\x1b[0m\r\n\x1b[39m\x1b[39mthird line",
			width: 100,
			lines: 2,
		},
		{
			msg:   "Multi-line string with color escape sequences and wrapping",
			input: "first line\nsecond line                                    \x1b[0m\x1b[0m\x1b[39m\x1b[39m\x1b[90m\x1b[90m | \x1b[0m\x1b[39m\x1b[0m\x1b[39m\x1b[0m\x1b[0m\x1b[39m\x1b[39mA separate column on the same line                \x1b[0m\x1b[0m\r\n\x1b[39m\x1b[39mthird line",
			width: 50,
			lines: 3,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.msg, func(t *testing.T) {
			got := numNewLinesInWidth(tt.input, tt.width)

			if got != tt.lines {
				t.Errorf("Returned line count was incorrect: got %d, want %d", got, tt.lines)
			}
		})
	}
}
