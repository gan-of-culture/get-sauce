package eahentai

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/test"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single Gallery",
			URL:  "https://eahentai.com/a/62293",
			Want: 1,
		}, {
			Name: "Collection",
			URL:  "https://eahentai.com/search?type=from&q=86&p=1",
			Want: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			URLs, err := parseURL(tt.URL)
			test.CheckError(t, err)
			if len(URLs) > tt.Want || len(URLs) == 0 {
				t.Errorf("Got: %v - Want: %v", len(URLs), tt.Want)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		Name string
		Args test.Args
	}{
		{
			Name: "Single Gallery",
			Args: test.Args{
				URL:   "https://eahentai.com/a/62293",
				Title: "(C107) [Fishspell (fnd)] Nee, Ii desho? (Yu-Gi-Oh! 5D's)",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.Args.URL)
			test.CheckError(t, err)
			test.Check(t, tt.Args, data[0])
		})
	}
}
