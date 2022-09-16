package miohentai

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
			URL:  "https://miohentai.com/enjo-kouhai-episode-2/",
			Want: 1,
		}, {
			Name: "Tag",
			URL:  "https://miohentai.com/tag/1080p/",
			Want: 20,
		}, {
			Name: "Image",
			URL:  "https://miohentai.com/image-library/my-favorite-sexy-lingerie/",
			Want: 1,
		},
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
				URL:     "https://miohentai.com/enjo-kouhai-episode-2/",
				Title:   "Enjo Kouhai – Episode 2",
				Quality: "",
				Size:    131533561,
			},
		},
		{
			Name: "Single Image",
			Args: test.Args{
				URL:     "https://miohentai.com/image-library/my-favorite-sexy-lingerie/",
				Title:   "my-favorite-sexy-lingerie",
				Quality: "",
				Size:    170503,
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
