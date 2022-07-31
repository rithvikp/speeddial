package term

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/pterm/pterm"
	"golang.org/x/exp/constraints"
)

// Canned errors for terminal interfaces.
var (
	ErrUserQuit = errors.New("the user manually exited the view")
)

// FormattedChunk defines a section of contiguous text.
type FormattedChunk struct {
	Start  int
	Length int
}

// FormattedContent is a set of text with certain parts marked for additional formatting. Content
// contains all text to be shown, with the Highlights slice specifying section that should be
// displayed differently.
type FormattedContent struct {
	Content    string
	Highlights []FormattedChunk
}

// ListItem represents an individual item in the list. It is made up of a list of content to display
// and an associated arbitrary piece of data that will be returned to the caller if the item is
// selected.
type ListItem[T any] struct {
	DisplayFields []FormattedContent
	Raw           T
}

// QueryableList abstracts a searchable corpus of data. The Search method will be repeatedly called as
// the query changes.
type QueryableList[T any] interface {
	Search(query string) []ListItem[T]
}

func min[T constraints.Ordered](a, b T) T {
	if a <= b {
		return a
	}
	return b
}

func max[T constraints.Ordered](a, b T) T {
	if a >= b {
		return a
	}
	return b
}

// List implements an interactive terminal list, printing the interface out to stderr and allowing
// the user to navigate and choose an option.
//
// The Vim bindings are currently limited to just list navigation.
func List[Payload any](list QueryableList[Payload], maxToDisplay int, vimNavigation bool) (Payload, error) {
	t, err := NewTty()
	defer func() {
		err := t.Stop()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to restore the terminal interface: %v", err)
		}
	}()

	var emptyPayload Payload
	if err != nil {
		return emptyPayload, fmt.Errorf("unable to initialize the terminal interface: %v", err)
	}

	query := ""
	items := list.Search(query)
	displayOffset := 0
	selected := 0
	normalMode := false
	output := ""

	moveCursorToStart := func() string {
		lines := max(0, t.NumLines(output)-1)
		return "\r" + vt100CursorUp(lines)
	}
	moveCursorToEndOfQuery := func(queryLen int) string {
		return fmt.Sprint(moveCursorToStart(), vt100CursorRight(len("> ")+queryLen))
	}

	listNavDown := func() {
		if selected < len(items)-1 {
			if selected-displayOffset >= maxToDisplay-1 {
				displayOffset++
			}
			selected++
		}
	}
	listNavUp := func() {
		if selected > 0 {
			if selected <= displayOffset {
				displayOffset--
			}
			selected--
		}
	}

	// Every iteration, first update the interface (print a single string with all the content
	// and relevant escape codes in order to have a smooth UI). Then, wait for a keystroke,
	// handle it appropriately, and repeat the entire process.
	for {
		// TODO: Refactor the terminal wrangling code to improve readability
		output = "\r"

		// Print the updated interface
		output += fmt.Sprint("> ", query, vt100ClearEOL(), "\r\n")

		tbl, err := generateList(t, items, displayOffset, maxToDisplay, selected)
		if err != nil {
			return emptyPayload, err
		}
		output += tbl

		// Wipe the rest of the screen downwards to remove old, trailing text
		output += vt100ClearEOS()
		output += moveCursorToEndOfQuery(len(query))

		fmt.Fprint(os.Stderr, output)

		// Handle keyboard events accordingly
		e, err := t.GetKeyboardEvent()
		if err != nil {
			return emptyPayload, fmt.Errorf("unable to process user keystroke: %v", err)
		}

		switch e.key {
		case KeyChar:
			if !normalMode {
				query += string(e.char)
				items = list.Search(query)
				selected = max(min(selected, len(items)-1), 0)
				displayOffset = max(min(displayOffset, len(items)-1-maxToDisplay), 0)
				break
			}

			if e.char == 'j' {
				listNavDown()
			} else if e.char == 'k' {
				listNavUp()
			} else if e.char == 'i' || e.char == 'a' {
				normalMode = false
			}

		case KeyEnter:
			if selected < 0 || selected >= len(items) {
				return emptyPayload, errors.New("unable to select an item")
			}
			// Wipe any added content added by this function
			fmt.Fprint(os.Stderr, moveCursorToStart(), vt100ClearEOS())

			return items[selected].Raw, nil

		case KeyCtrlC:
			// Wipe any added content added by this function
			fmt.Fprint(os.Stderr, "\r", vt100ClearEOS())
			return emptyPayload, ErrUserQuit

		case KeyDelete:
			if len(query) > 0 {
				query = query[:len(query)-1]
				items = list.Search(query)
				selected = max(min(selected, len(items)-1), 0)
				displayOffset = max(min(displayOffset, len(items)-1-maxToDisplay), 0)
			}

		case KeyUp:
			listNavUp()

		case KeyDown:
			listNavDown()

		case KeyEscape:
			if !vimNavigation {
				return emptyPayload, ErrUserQuit
			}

			normalMode = true
		}
	}
}

func generateList[T any](t *Tty, items []ListItem[T], displayOffset, maxToDisplay, selected int) (string, error) {
	if displayOffset < 0 || maxToDisplay < 0 {
		return "", fmt.Errorf("invalid display offset %d and/or range %d", displayOffset, maxToDisplay)
	} else if len(items) == 0 {
		return "", nil
	}

	endIndex := displayOffset + maxToDisplay
	if endIndex > len(items) {
		endIndex = len(items)
	}

	// Convert the list to a pterm table, bolding the selected row along the way
	var data pterm.TableData
	for i := displayOffset; i < endIndex; i++ {
		item := items[i]

		var formatted []string
		for _, elem := range item.DisplayFields {

			// Highlight matching text
			text, err := formatContent(elem.Content, elem.Highlights, func(s string) string {
				return pterm.Cyan(s)
			})
			if err != nil {
				return "", err
			}

			if i == selected {
				text = pterm.Bold.Sprint(text)
			}

			formatted = append(formatted, text)
		}

		data = append(data, formatted)
	}

	tbl, err := pterm.DefaultTable.WithData(data).Srender()
	if err != nil {
		return "", fmt.Errorf("unable to print the list: %v", err)
	}

	// Since the terminal is in raw mode, carriage returns are necessary, so add them in
	tbl = strings.ReplaceAll(tbl, "\n", vt100ClearEOL()+"\r\n")

	return tbl, nil
}

// This function should only be called on unformatted strings (unless the chunk indices take the
// formatting into account).
func formatContent(origContent string, chunks []FormattedChunk, f func(string) string) (string, error) {
	sort.SliceStable(chunks, func(i, j int) bool {
		return chunks[i].Start < chunks[j].Start
	})

	var fmtChunks []string
	j := 0
	for _, fc := range chunks {
		if fc.Start < 0 || fc.Length < 0 || fc.Start+fc.Length > len(origContent) {
			return "", fmt.Errorf("found an invalid format chunk with start %d and length %d", fc.Start, fc.Length)
		} else if fc.Start < j {
			return "", fmt.Errorf("found an overlapping format chunk with start %d", fc.Start)
		}

		fChunk := pterm.Cyan(origContent[fc.Start : fc.Start+fc.Length])
		fmtChunks = append(fmtChunks, origContent[j:fc.Start], fChunk)
		j = fc.Start + fc.Length
	}
	fmtChunks = append(fmtChunks, origContent[j:])

	return strings.Join(fmtChunks, ""), nil
}
