package main

import (
	"fmt"
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
			if err != nil {
				t.Errorf("Unexpected error: %s", err)
			}
			if page != tt.wantPage {
				t.Errorf("page; want: %d, got: %d", page, tt.wantPage)
			}
			if offset != tt.wantOffset {
				t.Errorf("offset; want: %d, got: %d", offset, tt.wantOffset)
			}
			if limit != tt.wantLimit {
				t.Errorf("limit; want: %d, got: %d", limit, tt.wantLimit)
			}
		})
	}
}

func TestLinkHeaders(t *testing.T) {
	fmt.Printf("%d\n", (37/5)+1)
}
