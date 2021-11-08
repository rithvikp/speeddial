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
func List(list QueryableList, maxToDisplay int) (interface{}, error) {
	// TODO: Figure out how to best handle errors/abstract this component
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

	moveCursorToStart := func(qlen, lines int) {
		// Move the cursor up to right below the prompt and left-align (for re-printing)
		fmt.Fprintf(os.Stderr, "\033[%dD", qlen+30)
		if lines > 0 {
			fmt.Fprintf(os.Stderr, "\033[%dA", lines)
		}
	}

	query := ""
	lines := 0
	items := list.Search(query)
	selected := 0
	for {
		addedLines, err := printList(t, items, maxToDisplay, selected)
		if err != nil {
			return nil, err
		}
		lines += addedLines

		fmt.Fprint(os.Stderr, "> ", query)

		// Wipe the rest of the screen, downwards
		fmt.Fprint(os.Stderr, "\033[0J")

		e, _ := t.GetKeyboardEvent()
		if e.key == KeyCtrlC || e.key == KeyEscape {
			return nil, ErrUserQuit
		}

		if e.key == KeyEnter {
			if selected < 0 || selected >= len(items) {
				return nil, errors.New("unable to select an item")
			}
			moveCursorToStart(len(query), lines)
			// Wipe any added content added by this function
			fmt.Fprint(os.Stderr, "\033[0J")

			return items[selected].Raw, nil
		}

		if e.key == KeyRune {
			query += string(e.char)
			items = list.Search(query)
			selected = max(min(selected, len(items)-1), 0)
		} else if e.key == KeyDelete && len(query) > 0 {
			query = query[:len(query)-1]
			items = list.Search(query)
			selected = max(min(selected, len(items)-1), 0)
		} else if e.key == KeyUp {
			if selected > 0 {
				selected -= 1
			}
		} else if e.key == KeyDown {
			if selected < len(items)-1 {
				selected += 1
			}
		}

		// TODO: Only re-render what has changed
		// TODO: Refactor all the terminal wrangling code

		moveCursorToStart(len(query), lines)
		lines = 0
	}
}

func printList(t *Tty, items []ListItem, maxToDisplay, selected int) (int, error) {
	if len(items) == 0 {
		return 0, nil
	}

	lines := 0

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
		return 0, fmt.Errorf("unable to print the list: %v", err)
	}

	tbl = strings.ReplaceAll(tbl, "\n", "\033[K\r\n")

	pterm.Fprint(os.Stderr, pterm.Sprint(tbl, "\033[K\r\n"))

	return lines, nil
}
