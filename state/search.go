package state

import (
	"strings"
)

// SearchRes defines the response from a command search.
type SearchRes struct {
	commands []*command
}

// Search searches all state in this container based on the given query.
func (c *Container) Search(query string) *SearchRes {
	var matched []*command

	for _, s := range c.states {
		matched = append(matched, s.search(query)...)
	}

	return &SearchRes{commands: matched}
}

// search searches the commands in this state to find any that match to the query. Currently,
// matching is purely based on "contains" operations.
func (s *state) search(query string) []*command {
	query = strings.ToLower(query)

	var matched []*command

	// TODO: Implement a robust text search scheme
	for _, c := range s.Commands {
		if strings.Contains(strings.ToLower(c.Invocation), query) {
			matched = append(matched, c)

		} else if strings.Contains(strings.ToLower(c.Description), query) {
			matched = append(matched, c)
		}
	}

	return matched
}

func (r *SearchRes) Get(i int) (string, string) {
	if i >= len(r.commands) {
		return "", ""
	}

	c := r.commands[i]
	return c.Invocation, c.Description
}

func (r *SearchRes) Size() int {
	return len(r.commands)
}
