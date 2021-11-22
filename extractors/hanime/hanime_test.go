package hanime

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single Episode hanime.io/",
			URL:  "https://hanime.io/hentai/torokase-orgasm-1/",
			Want: 1,
		}, {
			Name: "Overview hanime.io/",
			URL:  "https://hanime.io/genre/adventure/",
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
			Name: "Single Episode hanime.io/",
			URL:  "https://hanime.io/hentai/torokase-orgasm-1/",
			Want: 1,
		}, {
			Name: "Overview hanime.io/",
			URL:  "https://hanime.io/genre/adventure/",
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
