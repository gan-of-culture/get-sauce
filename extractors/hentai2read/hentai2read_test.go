package hentai2read

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single Gallery",
			URL:  "https://hentai2read.com/okitasan_to_kotasu_ecchi/#availableChapters",
			Want: 1,
		}, {
			Name: "Tag",
			URL:  "https://hentai2read.com/hentai-list/category/Romance/",
			Want: 48,
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
		URL  string
		Want int
	}{
		{
			Name: "Single Gallery",
			URL:  "https://hentai2read.com/elevenpm_miniature_garden/#availableChapters",
			Want: 1,
		}, /*{
			Name: "Tag",
			URL:  "https://hentai2read.com/hentai-list/category/Romance/",
			Want: 20,
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.URL)
			if err != nil {
				t.Error(err)
			}
			if len(data) > tt.Want || len(data) == 0 {
				t.Errorf("Got: %v - Want: %v", len(data), tt.Want)
			}
		})
	}
}
