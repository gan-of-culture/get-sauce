package hentai2w

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single Video",
			URL:  "https://hentai2w.com/video/youkoso-sukebe-elf-no-mori-e-episode-2-3693.html",
			Want: 1,
		}, {
			Name: "Category",
			URL:  "https://hentai2w.com/channels/125/magic/",
			Want: 40,
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
			Name: "Single Video",
			URL:  "https://hentai2w.com/video/youkoso-sukebe-elf-no-mori-e-episode-2-3693.html",
			Want: 1,
		}, /*{
			Name: "Category",
			URL:  "https://hentai2w.com/channels/125/magic/",
			Want: 40,
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
