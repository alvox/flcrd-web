package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

func TestPaging(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		wantPage   int
		wantOffset int
		wantLimit  int
	}{
		{
			name:       "Query in place",
			query:      "page=3&per_page=10",
			wantPage:   3,
			wantOffset: 20,
			wantLimit:  10,
		},
		{
			name:       "No query",
			query:      "",
			wantPage:   1,
			wantOffset: 0,
			wantLimit:  5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &http.Request{
				Method: "GET",
				URL: &url.URL{
					Scheme:   "https",
					Host:     "localhost",
					Path:     "public/decks",
					RawQuery: tt.query,
				},
			}
			page, offset, limit, err := extractPaging(r)
			require.Nil(t, err)
			require.Equal(t, tt.wantPage, page)
			require.Equal(t, tt.wantOffset, offset)
			require.Equal(t, tt.wantLimit, limit)
		})
	}
}

func TestLinkHeaders(t *testing.T) {
	fmt.Printf("%d\n", (37/5)+1)
}
