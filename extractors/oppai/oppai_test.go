package oppai

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
			Name: "Single Episode",
			URL:  "https://oppai.stream/watch.php?e=Kyonyuu-Elf-Oyako-Saimin-1&f=s",
			Want: 1,
		},
		// only extract single episodes for now since the search parameters are not saved in the URL
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			URLs := parseURL(tt.URL)
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
			Name: "Single Episode",
			Args: test.Args{
				URL:     "https://oppai.stream/watch.php?e=Kyonyuu-Elf-Oyako-Saimin-1&f=s",
				Title:   "Kyonyuu-Elf-Oyako-Saimin-1",
				Quality: "2160p",
				Size:    429522316,
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
