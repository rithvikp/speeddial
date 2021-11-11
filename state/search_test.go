package state

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMatch(t *testing.T) {
	tests := []struct {
		msg     string
		query   string
		src     string
		matches []matchedText
	}{
		{
			msg:   "no matches",
			query: "doc",
			src:   "source ./shell/setup.zsh",
		},
		{
			msg:   "exact equality match",
			query: "go build",
			src:   "go build",
			matches: []matchedText{
				{start: 0, length: 8},
			},
		},
		{
			msg:   "prefix equality match",
			query: "go build",
			src:   "go build github.com",
			matches: []matchedText{
				{start: 0, length: 8},
			},
		},
		{
			msg:   "suffix equality match",
			query: "github.com",
			src:   "go build github.com",
			matches: []matchedText{
				{start: 9, length: 10},
			},
		},
		{
			msg:   "interior equality match",
			query: "show pod",
			src:   "kubectl show pods",
			matches: []matchedText{
				{start: 8, length: 8},
			},
		},
		// TODO: The algorithm currently quite imperfect at minimizing the number of chunks.
		// It makes local decisions to minimize chunks, but, without look-ahead/look-behind,
		// this strategy cannot catch most cases.
		//{
		//msg:   "split equality prefix/suffix match",
		//query: "kube pods",
		//src:   "kubectl get pods",
		//matches: []matchedText{
		//{start: 0, length: 4},
		//{start: 11, length: 5},
		//},
		//},
		{
			msg:   "split equality interior match",
			query: "run hub",
			src:   "go run github.com",
			matches: []matchedText{
				{start: 3, length: 4},
				{start: 10, length: 3},
			},
		},
		{
			msg:   "take the first of multiple matches",
			query: "run",
			src:   "go run; go run; go run",
			matches: []matchedText{
				{start: 3, length: 3},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.msg, func(t *testing.T) {
			q := query{
				raw:     tt.query,
				cleaned: tt.query,
			}

			matches := match(&q, tt.src)
			if diff := cmp.Diff(matches, tt.matches, cmp.AllowUnexported(matchedText{})); diff != "" {
				t.Errorf("[]matchedText diff (-got, +want):\n%s", diff)
			}
		})
	}
}
