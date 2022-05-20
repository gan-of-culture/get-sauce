package hentaifoundry

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/test"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single Post",
			URL:  "https://www.hentai-foundry.com/pictures/user/LumiNyu/930609/Ty-Lee",
			Want: 1,
		}, {
			Name: "Overview User",
			URL:  "https://www.hentai-foundry.com/pictures/user/LumiNyu",
			Want: 50,
		}, {
			Name: "Overview Category",
			URL:  "https://www.hentai-foundry.com/categories/372/Anime-and-Manga/Chobits/pictures",
			Want: 10,
		},
	}
	config.Amount = 50
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
				URL:     "https://www.hentai-foundry.com/pictures/user/AyyaSAP/795835/Albedo",
				Title:   "AyyaSAP-795835-Albedo",
				Quality: "720x1008",
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
