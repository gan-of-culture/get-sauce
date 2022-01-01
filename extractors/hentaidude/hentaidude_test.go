package hentaidude

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
			URL:  "https://hentaidude.com/aisei-tenshi-love-mary-akusei-jutai-episode-1/",
			Want: 1,
		}, {
			Name: "Landing page",
			URL:  "https://hentaidude.com/",
			Want: 20,
		}, {
			Name: "Orderby",
			URL:  "https://hentaidude.com/?orderby=date",
			Want: 20,
		}, {
			Name: "Tags",
			URL:  "https://hentaidude.com/?orderby=date&tid=1472",
			Want: 20,
		}, {
			Name: "Different page",
			URL:  "https://hentaidude.com/page/3/?orderby=date&tid=1472",
			Want: 20,
		}, {
			Name: "3D Single Episode",
			URL:  "https://hentaidude.com/scarlet-nights-episode-1/",
			Want: 1,
		}, {
			Name: "3D Landing page",
			URL:  "https://hentaidude.com/tag/3d-hentai-0/",
			Want: 20,
		}, {
			Name: "3D Orderby",
			URL:  "https://hentaidude.com/tag/3d-hentai-0/?orderby=date",
			Want: 20,
		}, {
			Name: "3D Tags",
			URL:  "https://hentaidude.com/tag/3d-hentai-0/?orderby=date&tid=1472",
			Want: 20,
		}, {
			Name: "3D Different page",
			URL:  "https://hentaidude.com/tag/3d-hentai-0/page/2/?tid=1541",
			Want: 20,
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
			Name: "Single Episode",
			Args: test.Args{
				URL:     "https://hentaidude.com/yuutousei-ayaka-no-uraomote-episode-1/",
				Title:   "Yuutousei Ayaka No Uraomote - Episode 1",
				Quality: "",
				Size:    107194356,
			},
		},
		{
			Name: "Single 3D Episode",
			Args: test.Args{
				URL:     "https://hentaidude.com/scarlet-nights-episode-1/",
				Title:   "Scarlet Nights - Episode 1",
				Quality: "",
				Size:    436684542,
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
