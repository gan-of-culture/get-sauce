package hentaihavenred

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single Episode hentaihaven.red/",
			URL:  "https://hentaihaven.red/hentai/joshi-luck-episode-1/",
			Want: 1,
		}, {
			Name: "Overview hentaihaven.red/",
			URL:  "https://hentaihaven.red/ratings/",
			Want: 35,
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
			Name: "Single Episode hentaihaven.red",
			URL:  "https://hentaihaven.red/hentai/bitch-na-inane-sama-episode-4/",
			Want: 1,
		}, {
			Name: "[OLD] Single Episode hentaihaven.red",
			URL:  "https://hentaihaven.red/hentai/mako-chan-kaihatsu-nikki-episode-2/",
			Want: 1,
		}, {
			Name: "Overview hentaihaven.red",
			URL:  "https://hentaihaven.red/genre/2019-english/",
			Want: 4,
			//can be more videos at the time when I am adding this it was only 4 -> normally it is 30 per site but that would be too much for testing
		},
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
