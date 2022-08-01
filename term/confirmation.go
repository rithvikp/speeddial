package term

import (
	"fmt"
	"os"

	"github.com/rithvikp/speeddial/term/termui"
)

// Confirmation implements an interactive confirmation dialog. The corresponding message is printed
// out to stderr, with true being returned if the user confirms, false if not. If clearAfterUse is
// set, the confirmation dialog will be cleared before the function returns.
func Confirmation(msg string, clearAfterUse bool) (bool, error) {
	t, err := NewTty()
	defer func() {
		err := t.Stop()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to restore the terminal interface: %v", err)
		}
	}()
	if err != nil {
		return false, fmt.Errorf("unable to initialize the terminal interface: %v", err)
	}

	builder := &termui.Builder{}
	builder.SaveCursor()

	builder.WriteString(msg + " [y/n]")
	fmt.Fprint(os.Stderr, builder.Commit())

	e, err := t.GetKeyboardEvent()
	if err != nil {
		return false, fmt.Errorf("unable to process user keystroke: %v", err)
	}

	builder.ResetCursor().ClearToScreenEnd()
	fmt.Fprint(os.Stderr, builder.Commit())

	if e.key == KeyChar && e.char == 'y' {
		return true, nil
	} else if e.key == KeyCtrlC {
		return false, ErrUserQuit
	}
	return false, nil
}
