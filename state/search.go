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

// match takes a given query and determines whether it exists in src as a sequence of
// (potentially non-consecutive) substrings, returning the substrings if a match is found.
//
// A DP algorithm is used: it looks us successive prefixes of src and query, keeping track of
// the best matches so far.
func match(q *query, src string) []matchedText {
	if len(q.cleaned) == 0 || len(src) == 0 || len(q.cleaned) > len(src) {
		return nil
	}

	type state struct {
		valid   bool
		matches []matchedText
	}

	src = strings.TrimSpace(strings.ToLower(src))

	min := func(a, b int) int {
		if a <= b {
			return a
		}
		return b
	}

	dp := make([][]state, len(src))

	for i := 0; i < len(src); i++ {
		dp[i] = make([]state, len(q.cleaned))
		if src[i] == q.cleaned[0] {
			dp[i][0].valid = true
			dp[i][0].matches = append(dp[i][0].matches, matchedText{
				start:  i,
				length: 1,
			})
		}

		for j := 1; j < min(i+1, len(q.cleaned)); j++ {
			if src[i] == q.cleaned[j] && dp[i-1][j-1].valid {
				// Choose either to take a match that occurred before or add to a currently active
				// one depending on which minimizes the number of chunks.
				if !(dp[i-1][j].valid && len(dp[i-1][j-1].matches)+1 >= len(dp[i-1][j].matches)) {
					dp[i][j].valid = true
					dp[i][j].matches = append(dp[i][j].matches, dp[i-1][j-1].matches...)
					prev := &dp[i][j].matches[len(dp[i][j].matches)-1]

					if prev.start+prev.length == i {
						prev.length++
					} else {
						dp[i][j].matches = append(dp[i][j].matches, matchedText{
							start:  i,
							length: 1,
						})
					}
					break
				}
			}

			if dp[i-1][j].valid {
				dp[i][j].valid = true
				dp[i][j].matches = append(dp[i][j].matches, dp[i-1][j].matches...)
			}
		}
	}

	return dp[len(src)-1][len(q.cleaned)-1].matches
}
