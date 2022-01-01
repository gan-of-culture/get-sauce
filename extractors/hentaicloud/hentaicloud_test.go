package hentaicloud

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
			URL:  "https://www.hentaicloud.com/video/3366/rikujoubu-joshi-wa-ore-no-nama-onaho-the-animation/episode2/english",
			Want: 1,
		}, {
			Name: "Group",
			URL:  "https://www.hentaicloud.com/videos/oppai",
			Want: 23,
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
				URL:     "https://www.hentaicloud.com/video/750/boy-meets-harem-the-animation/episode1",
				Title:   "Boy Meets Harem The Animation",
				Quality: "720p",
				Size:    244322800,
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
