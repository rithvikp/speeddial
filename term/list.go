package term

import (
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"
)

type QueryableList interface {
	Search(query string) [][]string
}

func List(list QueryableList) {
	// TODO: Figure out how to best handle errors/abstract this component
	t, err := NewTty()
	defer t.Stop()
	if err != nil {
		fmt.Println("Unable to initialize the terminal interface")
		return
	}

	min := func(a, b int) int {
		if a <= b {
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
		lines += printList(t, items, selected)
		fmt.Fprint(os.Stderr, "> ", query)
		// Wipe the rest of the screen, downwards
		fmt.Fprint(os.Stderr, "\033[0J")

		e, _ := t.GetKeyboardEvent()
		if e.key == KeyCtrlC || e.key == KeyEscape {
			return
		}

		if e.key == KeyEnter {
			if selected < 0 || selected >= len(items) {
				fmt.Println("Unable to select an item")
				return
			}
			moveCursorToStart(len(query), lines)
			// Wipe any added content added by this function
			fmt.Fprint(os.Stderr, "\033[0J")
			fmt.Print(items[selected][0], "\r\n")
			return
		}

		if e.key == KeyRune {
			query += string(e.char)
			items = list.Search(query)
			selected = min(selected, len(items)-1)
		} else if e.key == KeyDelete && len(query) > 0 {
			query = query[:len(query)-1]
			items = list.Search(query)
			selected = min(selected, len(items)-1)
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

func printList(t *Tty, items [][]string, selected int) int {
	if len(items) == 0 {
		return 0
	}

	lines := 0

	var data pterm.TableData
	for i, item := range items {
		lines++
		if i != selected {
			data = append(data, item)
			continue
		}

		var formatted []string
		for _, elem := range item {
			formatted = append(formatted, pterm.Bold.Sprint(elem))
		}
		data = append(data, formatted)
	}

	tbl, err := pterm.DefaultTable.WithData(data).Srender()
	if err != nil {
		fmt.Println("Unable to print the list")
		return 1
	}

	tbl = strings.ReplaceAll(tbl, "\n", "\033[K\r\n")

	pterm.Fprint(os.Stderr, pterm.Sprint(tbl, "\033[K\r\n"))

	return lines
}
