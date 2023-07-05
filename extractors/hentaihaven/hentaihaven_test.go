package hentaihaven

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
			Name: "Single Episode hentaihaven.com",
			URL:  "https://hentaihaven.com/video/soshite-watashi-wa-sensei-ni/episode-1/",
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
			Name: "Single Episode hentaihaven.com",
			Args: test.Args{
				URL:     "https://hentaihaven.com/video/usamimi-bouken-tan-sekuhara-shinagara-sekai-o-sukue/episode-3/",
				Title:   "Usamimi Bouken-tan: Sekuhara Shinagara Sekai o Sukue - Episode 3",
				Quality: "1920x1080",
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
