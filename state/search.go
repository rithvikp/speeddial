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

// Searcher returns a Searcher over the container.
func (c *Container) Searcher() *Searcher {
	return &Searcher{c: c}
}

type query struct {
	raw      string
	cleaned  string
	unigrams []string
}

type matchedText struct {
	start  int
	length int
}

type matchedCommand struct {
	c           *Command
	invMatches  []matchedText
	descMatches []matchedText
}

// Search searches all state in this container based on the given query.
func (s *Searcher) Search(rawQuery string) []term.ListItem[*Command] {
	var matched []term.ListItem[*Command]
	q := parseQuery(rawQuery)

	for _, s := range s.c.states {
		for _, m := range s.search(q) {
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

			li := term.ListItem[*Command]{
				DisplayFields: []term.FormattedContent{inv, desc},
				Raw:           m.c,
			}

			matched = append(matched, li)
		}
	}

	return matched
}

func parseQuery(raw string) *query {
	q := query{
		raw:     raw,
		cleaned: strings.TrimSpace(strings.ToLower(raw)),
	}
	q.unigrams = strings.Fields(q.cleaned)
	return &q
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

		if q.cleaned == "" {
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

// match takes a given query and determines whether it exists in src as a sequence of
// (potentially non-consecutive) substrings, returning the substrings if a match is found.
//
// A DP algorithm is used: it looks at successive prefixes of src and query, keeping track of
// all matches so far. At the end, the full match (ie. a match where the entire query appears
// in some form in src) with the fewest chunks (and then the earliest chunk as a tiebreaker)
// is returned.
//
// TODO: Minimize the number of tokens that are split across chunks.
func match(q *query, src string) []matchedText {
	if q.cleaned == "" || src == "" || len(q.cleaned) > len(src) {
		return nil
	}

	type state struct {
		valid   bool
		matches [][]matchedText
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
			dp[i][0].matches = append(dp[i][0].matches, []matchedText{{
				start:  i,
				length: 1,
			}})
		} else if i > 0 && dp[i-1][0].valid {
			// If an element of the first column is valid, it will have exactly one match list
			dp[i][0].valid = true
			dp[i][0].matches = append(dp[i][0].matches, []matchedText{})
			dp[i][0].matches[0] = append(dp[i][0].matches[0], dp[i-1][0].matches[0]...)
		}

		for j := 1; j < min(i+1, len(q.cleaned)); j++ {
			if src[i] == q.cleaned[j] && dp[i-1][j-1].valid {
				for k := 0; k < len(dp[i-1][j-1].matches); k++ {
					dp[i][j].valid = true
					dp[i][j].matches = append(dp[i][j].matches, []matchedText{})
					dp[i][j].matches[k] = append(dp[i][j].matches[k], dp[i-1][j-1].matches[k]...)

					prev := &dp[i][j].matches[k][len(dp[i][j].matches[k])-1]
					if prev.start+prev.length == i {
						prev.length++
					} else {
						dp[i][j].matches[k] = append(dp[i][j].matches[k], matchedText{
							start:  i,
							length: 1,
						})
					}
				}
			}

			if dp[i-1][j].valid {
				offset := len(dp[i][j].matches)
				for k := 0; k < len(dp[i-1][j].matches); k++ {
					dp[i][j].valid = true
					dp[i][j].matches = append(dp[i][j].matches, []matchedText{})
					dp[i][j].matches[offset+k] = append(dp[i][j].matches[offset+k], dp[i-1][j].matches[k]...)
				}
			}
		}
	}

	var bestMatch []matchedText
	matches := dp[len(src)-1][len(q.cleaned)-1].matches
	for _, m := range matches {
		if len(bestMatch) == 0 || len(m) < len(bestMatch) {
			bestMatch = m
		} else if len(m) == len(bestMatch) {
			// As a precondition, all match lists should be non-empty, and chunks should be ordered
			// in increasing order of index.
			i := 0
			for i < len(m) && m[i].start == bestMatch[i].start {
				i++
			}

			if i < len(m) && m[i].start < bestMatch[i].start {
				bestMatch = m
			}
		}
	}

	return bestMatch
}
