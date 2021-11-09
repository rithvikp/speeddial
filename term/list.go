package term

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"
)

var (
	ErrUserQuit = errors.New("the user manually exited the view")
)

type ListItem struct {
	DisplayFields []string
	Raw           interface{}
}

type QueryableList interface {
	Search(query string) []ListItem
}

// TODO: I really don't like returning an interface{} and then casting it back to the correct type
// later. Need to brainstorm better ways to keep this component generalizable while cleaning up the
// interface

// The Vim bindings are currently limited to just list navigation.
func List(list QueryableList, maxToDisplay int, vimNavigation bool) (interface{}, error) {
	t, err := NewTty()
	defer t.Stop()
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the terminal interface: %v", err)
	}

	min := func(a, b int) int {
		if a <= b {
			return a
		}
		return b
	}

	max := func(a, b int) int {
		if a >= b {
			return a
		}
		return b
	}

	moveCursorToStart := func(lines int) string {
		output := "\r"
		if lines > 0 {
			output += vt100CursorUp(lines)
		}
		return output
	}

	moveCursorToEndOfQuery := func(qlen, lines int) string {
		return fmt.Sprint(moveCursorToStart(lines), vt100CursorRight(len("> ")+qlen))
	}

	query := ""
	lines := 0
	items := list.Search(query)
	selected := 0
	normalMode := false

	listNavDown := func() {
		if selected < min(len(items)-1, maxToDisplay-1) {
			selected += 1
		}
	}

	listNavUp := func() {
		if selected > 0 {
			selected -= 1
		}
	}

	// Every iteration, first update the interface (print a single string with all the content
	// and relevant escape codes in order to have a smooth UI). Then, wait for a keystroke,
	// handle it appropriately, and repeat the entire process.
	for {
		output := ""

		// TODO: Refactor the terminal wrangling code to improve readability
		output += moveCursorToStart(lines)

		// Print the updated interface
		output += fmt.Sprint("> ", query, vt100ClearEOL(), "\r\n")
		lines++

		tbl, addedLines, err := printList(t, items, maxToDisplay, selected)
		if err != nil {
			return nil, err
		}
		output += tbl
		lines += addedLines

		// Wipe the rest of the screen downwards to remove old, trailing text
		output += vt100ClearEOS()
		output += moveCursorToEndOfQuery(len(query), lines)
		lines = 0

		fmt.Fprint(os.Stderr, output)

		// Handle keyboard events accordingly
		e, err := t.GetKeyboardEvent()
		if err != nil {
			return nil, fmt.Errorf("unable to process user keystroke: %v", err)
		}

		switch e.key {
		case KeyRune:
			if !normalMode {
				query += string(e.char)
				items = list.Search(query)
				selected = max(min(selected, len(items)-1), 0)
				break
			}

			if e.char == 'j' {
				listNavDown()
			} else if e.char == 'k' {
				listNavUp()
			} else if e.char == 'i' {
				normalMode = false
			}

		case KeyEnter:
			if selected < 0 || selected >= len(items) {
				return nil, errors.New("unable to select an item")
			}
			// Wipe any added content added by this function
			fmt.Fprint(os.Stderr, moveCursorToStart(lines), vt100ClearEOS())

			return items[selected].Raw, nil

		case KeyCtrlC:
			// Wipe any added content added by this function
			fmt.Fprint(os.Stderr, moveCursorToStart(lines), vt100ClearEOS())
			return nil, ErrUserQuit

		case KeyDelete:
			if len(query) > 0 {
				query = query[:len(query)-1]
				items = list.Search(query)
				selected = max(min(selected, len(items)-1), 0)
			}

		case KeyUp:
			listNavUp()

		case KeyDown:
			listNavDown()

		case KeyEscape:
			if !vimNavigation {
				return nil, ErrUserQuit
			}

			normalMode = true
		}
	}
}

func printList(t *Tty, items []ListItem, maxToDisplay, selected int) (string, int, error) {
	if len(items) == 0 {
		return "", 0, nil
	}

	lines := 0

	// Convert the list to a pterm table, bolding the selected row along the way
	var data pterm.TableData
	for i, item := range items {
		if i >= maxToDisplay {
			break
		}
		lines++
		if i != selected {
			data = append(data, item.DisplayFields)
			continue
		}

		var formatted []string
		for _, elem := range item.DisplayFields {
			formatted = append(formatted, pterm.Bold.Sprint(elem))
		}
		data = append(data, formatted)
	}

	tbl, err := pterm.DefaultTable.WithData(data).Srender()
	if err != nil {
		return "", 0, fmt.Errorf("unable to print the list: %v", err)
	}

	// Since the terminal is in raw mode, carriage returns are necessary, so add them in
	tbl = strings.ReplaceAll(tbl, "\n", vt100ClearEOL()+"\r\n")

	output := pterm.Sprint(tbl, vt100ClearEOL()+"\r\n")

	return output, lines, nil
}
