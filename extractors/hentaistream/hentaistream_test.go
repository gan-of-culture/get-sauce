package hentaistream

import "testing"

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
		URL  string
		Want int
	}{
		{
			Name: "Series Extraction 4k",
			URL:  "https://hentaistream.moe/anime/overflow/",
			Want: 8,
		}, {
			Name: "Single Episode 4k",
			URL:  "https://hentaistream.moe/515/overflow-1/",
			Want: 1,
		}, {
			Name: "Single Episode",
			URL:  "https://hentaistream.moe/593/kateikyoushi-no-oneesan-2-the-animation-1/",
			Want: 1,
		}, {
			Name: "Series",
			URL:  "https://hentaistream.moe/anime/kateikyoushi-no-oneesan-2-the-animation/",
			Want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.URL)
			if err != nil {
				t.Error(err)
			}
			if len(data) > tt.Want {
				t.Errorf("Got: %v - Want: %v", len(data), tt.Want)
			}
		})
	}
}
