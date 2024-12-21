package hentaiworld

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
			URL:  "https://hentaiworld.tv/hentai-videos/ijirare-fukushuu-saimin-episode-2/",
			Want: 1,
		}, {
			Name: "All episodes page",
			URL:  "https://hentaiworld.tv/all-episodes/page/2/",
			Want: 30,
		}, {
			Name: "Uncensored page",
			URL:  "https://hentaiworld.tv/uncensored/",
			Want: 30,
		}, {
			Name: "3d page",
			URL:  "https://hentaiworld.tv/3d/",
			Want: 60,
		}, {
			Name: "tag page",
			URL:  "https://hentaiworld.tv/hentai-videos/tag/anal/",
			Want: 30,
		}, {
			Name: "3d post page",
			URL:  "https://hentaiworld.tv/hentai-videos/3d/final-fantasy-tifa-7/",
			Want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			URLs := parseURL(tt.URL)
			if len(URLs) > tt.Want {
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
			Name: "Single Extraction",
			Args: test.Args{
				URL:     "https://hentaiworld.tv/hentai-videos/yuutousei-ayaka-no-uraomote-episode-1/",
				Title:   "Yuutousei Ayaka no Uraomote - Episode 1",
				Quality: "",
				Size:    80671065,
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
