package term

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/term"
)

// Tty represents a raw terminal interface.
type Tty struct {
	oldState *term.State
}

// Key represents keyboard keys.
type Key int

// Define various keys that are currently tracked by the application. These are mostly special
// keys/key combos. KeyChar represents an actual character. The definitions are split up into
// multiple blocks as only a few keys/key combos are currently tracked.
const (
	KeyChar Key = iota
	KeyCtrlA
	KeyCtrlB
	KeyCtrlC
)

// Define additional keys.
const (
	KeyUp Key = iota + 512
	KeyDown
)

// Define even more keys.
const (
	KeyEnter  Key = 13
	KeyEscape Key = 27
	KeyDelete Key = 127
)

// Event represents a keyboard event. If the Key is KeyChar, the char field should be checked for
// the specific character that was pressed.
type Event struct {
	key  Key
	char rune
}

// NewTty creates a new Tty. It has a side-effect of switching the current terminal to raw mode.
func NewTty() (*Tty, error) {
	t := Tty{}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}
	t.oldState = oldState

	return &t, err
}

// GetKeyboardEvent blocks until there is a keyboard event, and then returns it.
func (t *Tty) GetKeyboardEvent() (*Event, error) {
	buf := make([]byte, 5)
	n, err := os.Stdin.Read(buf)
	if err != nil {
		return nil, err
	} else if n == 0 {
		return nil, errors.New("unable to read any characters from tty")
	}

	buf = buf[:n]

	e := Event{}
	// Handle regular characters
	if n == 1 && buf[0] >= ' ' && buf[0] <= '~' {
		e.key = KeyChar
		e.char = rune(buf[0])
		return &e, nil
	}

	// TODO: Improve key matching (match more characters, match arrows in a better manner etc.)
	switch int(buf[0]) {
	case int(KeyCtrlC):
		e.key = KeyCtrlC
	case int(KeyEnter):
		e.key = KeyEnter
	case int(KeyEscape):
		e.key = KeyEscape
	case int(KeyDelete):
		e.key = KeyDelete
	}

	if len(buf) == 3 {
		if buf[0] == 27 && buf[1] == 91 && buf[2] == 65 {
			e.key = KeyUp
		} else if buf[0] == 27 && buf[1] == 91 && buf[2] == 66 {
			e.key = KeyDown
		}
	}

	return &e, nil
}

// Stop restore the current terminal to its previous state. It should be called after the caller
// is done using the Tty.
func (t *Tty) Stop() error {
	return term.Restore(int(os.Stdin.Fd()), t.oldState)
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
