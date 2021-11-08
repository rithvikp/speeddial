package term

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/term"
)

type Tty struct {
	oldState *term.State
}

type Key int

const (
	KeyRune Key = iota
	KeyCtrlA
	KeyCtrlB
	KeyCtrlC
)

const (
	KeyEnter  Key = 13
	KeyEscape Key = 27
	KeyDelete Key = 127
)

const (
	KeyUp Key = iota + 512
	KeyDown
)

type Event struct {
	key  Key
	char rune
}

func NewTty() (*Tty, error) {
	t := Tty{}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}
	t.oldState = oldState

	return &t, err
}

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
		e.key = KeyRune
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

func vt100CursorLeft(cols int) string {
	return fmt.Sprintf("\033[%dD", cols)
}
