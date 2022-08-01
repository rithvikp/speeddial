package termui

import (
	"fmt"
	"strings"

	"github.com/pterm/pterm"
	"golang.org/x/term"
)

// Builder builds inline UIs for terminal CLIs.
type Builder struct {
	*builder
	// The file descriptor which corresponds to the input from the terminal. This defaults to stdin.
	fd int
}

type builder struct {
	content        strings.Builder
	linesSinceSave int
	saveCol        int
}

func (b *Builder) init() {
	if b.builder == nil {
		b.builder = &builder{}
	}
}

func (b *Builder) width() int {
	b.init()

	width, _, err := term.GetSize(b.fd)
	if err != nil {
		panic(err)
	}
	return width
}

// Commit dumps the current contents to a strings and resets the internal buffer.
func (b *Builder) Commit() string {
	b.init()
	output := b.content.String()
	b.content.Reset()
	return output
}

// NextLine adds a new line to the buffer, clearing the remainder of the current line before
// doing so.
func (b *Builder) NextLine() *Builder {
	b.init()

	b.linesSinceSave++
	b.content.WriteString(vt100ClearEOL() + "\n\r")
	return b
}

// WriteStringAndReformat performs the same functionality as WriteString but it also
// reformats the provided string, adding carriage returns and clearing to the end of lines.
//
// This method is useful when a multi-line output is constructed using another tool but needs
// to be converted to work in a raw terminal.
func (b *Builder) WriteStringAndReformat(text string) *Builder {
	b.init()

	b.linesSinceSave += numNewLinesInWidth(text, b.width())

	// Since the terminal is in raw mode, carriage returns as well as clearing
	// the remainder of the line are necessary. This must be performed after counting
	// the number of new lines as numNewLinesInWidth does not ignore vt100 codes
	text = strings.ReplaceAll(text, "\n", vt100ClearEOL()+"\r\n")
	b.content.WriteString(text)
	return b
}

// WriteString adds the given string to the buffer.
func (b *Builder) WriteString(text string) *Builder {
	b.init()

	b.linesSinceSave += numNewLinesInWidth(text, b.width())
	b.content.WriteString(text)
	return b
}

// ClearToLineEnd clears any subsequent text on the current line.
func (b *Builder) ClearToLineEnd() *Builder {
	b.init()

	b.content.WriteString(vt100ClearEOL())
	return b
}

// ClearToScreenEnd clears any text from the current cursor position to the bottom of the terminal.
func (b *Builder) ClearToScreenEnd() *Builder {
	b.init()

	b.content.WriteString(vt100ClearEOS())
	return b
}

// SaveCursor saves the current position of the cursor for future resets. The vertical position
// is not a row number but an actual line (so may move if the window is scrolled).
func (b *Builder) SaveCursor() *Builder {
	b.init()

	b.linesSinceSave = 0
	b.saveCol = 0 // FIXME: Fetch col
	return b
}

// ResetCursor moves the cursor back to the last savepoint.
func (b *Builder) ResetCursor() *Builder {
	b.init()

	lines := b.linesSinceSave
	b.linesSinceSave = 0

	actions := make([]CursorAction, 0, 3)
	actions = append(actions, CursorLineStart())
	if lines > 0 {
		actions = append(actions, CursorUp(lines))
	} else if b.saveCol > 0 {
		actions = append(actions, CursorRight(b.saveCol))
	}

	return b.MoveCursor(actions...)
}

// MoveCursor moves the cursor as specified in the provided actions. This can include manual
// movements like moving up by two rows or commands like "move to the beginning of a line."
func (b *Builder) MoveCursor(actions ...CursorAction) *Builder {
	b.init()

	for _, a := range actions {
		b.content.WriteString(string(a))
	}
	return b
}

// CursorAction represents some movement that the cursor should perform.
type CursorAction string

// CursorRight translates the cursor right by the given number of columns.
func CursorRight(cols int) CursorAction {
	return CursorAction(vt100CursorRight(cols))
}

// CursorUp translates the cursor up by the given number of columns.
func CursorUp(rows int) CursorAction {
	return CursorAction(vt100CursorUp(rows))
}

// CursorLineStart moves the cursor to the beginning of the current line.
func CursorLineStart() CursorAction {
	return CursorAction("\r")
}

// numNewLinesInWidth returns the number of new lines (after wrapping) taken up by the
// given string when shown on a screen with the given width. There are some
// significant known limitations to the implementation (including support
// for vt100 escape codes, tabs etc.).
//
// Keep track of the number of characters in the current line, adding a new line
// when the max width is reached.
func numNewLinesInWidth(s string, width int) int {
	lines := 0
	charsInLine := 0
	s = pterm.RemoveColorFromString(s)

	for _, c := range s {
		if c == '\n' || charsInLine > width {
			charsInLine = 0
			lines++
		} else if c != '\r' {
			charsInLine++
		}
	}
	return lines
}

// Define functions for manipulating the screen of a VT100 terminal. These are defined as functinos
// and not constants as some have dynamically adjustable content (such as the number of lines to
// remove etc.).

func vt100ClearEOS() string {
	return "\033[0J"
}

func vt100ClearEOL() string {
	return "\033[K"
}

func vt100CursorUp(lines int) string {
	return fmt.Sprintf("\033[%dA", lines)
}

func vt100CursorRight(cols int) string {
	return fmt.Sprintf("\033[%dC", cols)
}

var _ = vt100CursorMove
var _ = vt100SaveCursorPos
var _ = vt100ResetCursorToSavedPos

func vt100CursorMove(row, col int) string {
	return fmt.Sprintf("\033[%d;%dH", row, col)
}

func vt100SaveCursorPos() string {
	return "\0337"
}

func vt100ResetCursorToSavedPos() string {
	return "\0338"
}
