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

type query struct {
	raw      string
	cleaned  string
	unigrams []string
}

type matchedCommand struct {
	c           *Command
	invMatches  []matchedText
	descMatches []matchedText
}

// Search searches all state in this container based on the given query.
func (s *Searcher) Search(rawQuery string) []term.ListItem {
	var matched []term.ListItem

	q := query{
		raw:     rawQuery,
		cleaned: strings.TrimSpace(strings.ToLower(rawQuery)),
	}
	q.unigrams = strings.Fields(q.cleaned)

	for _, s := range s.c.states {
		for _, m := range s.search(&q) {
			inv := term.FormattedContent{
				Content: m.c.Invocation,
			}

			for _, mt := range m.invMatches {
				inv.Highlights = append(inv.Highlights, term.FormattedChunk{
					Start:  mt.start,
					Length: mt.length,
				})
			}

			desc := term.FormattedContent{
				Content: m.c.Description,
			}

			for _, mt := range m.descMatches {
				desc.Highlights = append(desc.Highlights, term.FormattedChunk{
					Start:  mt.start,
					Length: mt.length,
				})
			}

			li := term.ListItem{
				DisplayFields: []term.FormattedContent{inv, desc},
				Raw:           m.c,
			}

			matched = append(matched, li)
		}
	}

	return matched
}

// search searches the commands in this state to find any that match to the query. Currently,
// matching is purely based on "contains" operations.
func (s *state) search(q *query) []matchedCommand {
	var matched []matchedCommand

	for _, c := range s.Commands {
		c.state = s
		mc := matchedCommand{
			c: c,
		}

		if len(q.cleaned) == 0 {
			matched = append(matched, mc)
			continue
		}

		mc.invMatches = match(q, c.Invocation)
		mc.descMatches = match(q, c.Description)
		if len(mc.invMatches) > 0 || len(mc.descMatches) > 0 {
			matched = append(matched, mc)
		}
	}

	return matched
}

type matchedText struct {
	start  int
	length int
}

func match(q *query, src string) []matchedText {
	if len(q.cleaned) == 0 {
		return nil
	}

	src = strings.TrimSpace(strings.ToLower(src))

	var matches []matchedText

	i := 0
	for i < len(src) {
		j := strings.Index(src[i:], q.cleaned)
		if j == -1 {
			break
		}

		mt := matchedText{
			start:  i + j,
			length: len(q.cleaned),
		}
		matches = append(matches, mt)

		i += j + len(q.cleaned)
	}

	return matches
}
