package term

import (
	"fmt"
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
		fmt.Printf("\033[%dD", qlen+30)
		if lines > 0 {
			fmt.Printf("\033[%dA", lines)
		}
	}

	query := ""
	lines := 0
	items := list.Search(query)
	selected := 0
	for {
		lines += printList(t, items, selected)
		fmt.Print("> ", query)
		// Wipe the rest of the screen, downwards
		fmt.Print("\033[0J")

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
			fmt.Print("\033[0J")
			fmt.Print(items[selected][0], "\033[K\r\n")
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
	//tabw := new(tabwriter.Writer)
	//tabw.Init(os.Stdout, 0, 2, 1, '\t', 0)
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

	pterm.Print(pterm.Sprint(tbl, "\033[K\r\n"))

	return lines

	//for i, item := range items {
	//for j, text := range item {
	//if i == selected {
	//fmt.Fprint(tabw, pterm.Bold.Sprint(text))
	//} else {
	//fmt.Fprint(tabw, text, pterm.Bold.Sprint())
	//}

	//if j < len(item)-1 {
	//fmt.Fprint(tabw, "\t")
	//} else {
	//fmt.Fprint(tabw, "\033[K")
	//}
	//}

	//// TODO: Cleanup the \n\r abstractions
	//fmt.Fprint(tabw, "\n\r")
	//lines += 1
	//}

	//if err := tabw.Flush(); err != nil {
	//fmt.Println("Unable to print the list")
	//return 0
	//}

	//return lines
}
