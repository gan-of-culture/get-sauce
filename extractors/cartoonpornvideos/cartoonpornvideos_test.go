package cartoonpornvideos

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
			URL:  "https://www.cartoonpornvideos.com/video/high-school-girl-who-was-groped-2-busty-hentai-teen-wants-to-fuck-in-the-internet-cafe-2Faxj2LF9fw.html",
			Want: 1,
		}, {
			Name: "Overview",
			URL:  "https://www.cartoonpornvideos.com/tags/video/curvy",
			Want: 34,
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
				URL:     "https://www.cartoonpornvideos.com/video/high-school-girl-who-was-groped-2-busty-hentai-teen-wants-to-fuck-in-the-internet-cafe-2Faxj2LF9fw.html",
				Title:   "High School Girl Who Was Groped 2 - Busty hentai teen wants to fuck in the internet cafe",
				Quality: "1920x1080",
				Size:    517783960,
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
