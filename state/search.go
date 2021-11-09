package state

import (
	"strings"

	"github.com/rithvikp/speeddial/term"
)

// Searcher provides a searchable view over all commands. It conforms to the
// term.QueryableList interface
type Searcher struct {
	c *Container
}

func (c *Container) Searcher() *Searcher {
	return &Searcher{c: c}
}

// Search searches all state in this container based on the given query.
func (s *Searcher) Search(query string) []term.ListItem {
	var matched []term.ListItem

	for _, s := range s.c.states {
		for _, m := range s.search(query) {
			li := term.ListItem{
				DisplayFields: []string{m.Invocation, m.Description},
				Raw:           m,
			}

			matched = append(matched, li)
		}
	}

	return matched
}

// search searches the commands in this state to find any that match to the query. Currently,
// matching is purely based on "contains" operations.
func (s *state) search(query string) []*Command {
	query = strings.ToLower(query)

	var matched []*Command

	// TODO: Implement a robust text search scheme
	for _, c := range s.Commands {
		c.state = s

		if strings.Contains(strings.ToLower(c.Invocation), query) {
			matched = append(matched, c)

		} else if strings.Contains(strings.ToLower(c.Description), query) {
			matched = append(matched, c)
		}
	}

	return matched
}
