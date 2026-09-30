package mhtdoujin

import (
	"net/url"
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
			Name: "Single Gallery HentaiEnvy",
			URL:  "https://hentaienvy.com/gallery/1043551/",
			Want: 1,
		}, {
			Name: "Tag HentaiEnvy",
			URL:  "https://hentaienvy.com/parody/azur-lane/",
			Want: 28,
		}, {
			Name: "Single Gallery HentaiZap",
			URL:  "https://hentaizap.com/gallery/843645/",
			Want: 1,
		}, {
			Name: "Tag HentaiZap",
			URL:  "https://hentaizap.com/tag/ahegao/",
			Want: 24,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			u, err := url.Parse(tt.URL)
			test.CheckError(t, err)

			URLs, err := parseURL(u)
			test.CheckError(t, err)
			if len(URLs) != tt.Want {
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
			Name: "Single Gallery HentaiEnvy",
			Args: test.Args{
				URL:   "https://hentaienvy.com/gallery/273160/",
				Title: "Makura Eigyou de Oshioki yo! ~Ano Sailor Senshi ga Makura Eigyou Halation~",
			},
		},
		{
			Name: "Single Gallery HentaiZap",
			Args: test.Args{
				URL:   "https://hentaizap.com/gallery/843645/",
				Title: "SUMMER FOX HUNTING",
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
