package damn

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single Episode",
			URL:  "https://www.damn.stream/watch/elf-hime-nina-episode-3",
			Want: 1,
		}, {
			Name: "Series",
			URL:  "https://www.damn.stream/hentai/ane-koi-suki-kirai-daisuki",
			Want: 2,
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
			Name: "Single Episode",
			URL:  "https://www.damn.stream/watch/elf-hime-nina-episode-3",
			Want: 1,
		}, {
			Name: "Series",
			URL:  "https://www.damn.stream/hentai/ane-koi-suki-kirai-daisuki",
			Want: 2,
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
