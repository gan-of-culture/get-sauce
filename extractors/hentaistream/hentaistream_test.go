package hentaistream

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
			URL:  "https://hentaistream.moe/593/kateikyoushi-no-oneesan-2-the-animation-1/",
			Want: 1,
		}, {
			Name: "Series",
			URL:  "https://hentaistream.moe/anime/kateikyoushi-no-oneesan-2-the-animation/",
			Want: 2,
		}, {
			Name: "Single Episode 4k",
			URL:  "https://hentaistream.moe/515/overflow-1/",
			Want: 1,
		}, {
			Name: "Series 4k",
			URL:  "https://hentaistream.moe/anime/overflow/",
			Want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			URLs := parseURL(tt.URL)
			if len(URLs) > tt.Want && tt.Want != 0 {
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
			Name: "Single Episode 4k",
			Args: test.Args{
				URL:     "https://hentaistream.moe/515/overflow-1/",
				Title:   "Overflow 1",
				Quality: "av1.2160p.webm",
				Size:    556,
			},
		},
		{
			Name: "Single Episode 4k",
			Args: test.Args{
				URL:     "https://hentaistream.moe/593/kateikyoushi-no-oneesan-2-the-animation-1/",
				Title:   "Kateikyoushi no Oneesan 2 The Animation 1",
				Quality: "av1.1080p.webm",
				Size:    555,
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
