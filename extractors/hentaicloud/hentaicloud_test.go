package hentaicloud

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single Episode",
			URL:  "https://www.hentaicloud.com/video/3366/rikujoubu-joshi-wa-ore-no-nama-onaho-the-animation/episode2/english",
			Want: 1,
		}, {
			Name: "Group",
			URL:  "https://www.hentaicloud.com/videos/oppai",
			Want: 23,
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
			URL:  "https://www.hentaicloud.com/video/750/boy-meets-harem-the-animation/episode1",
			Want: 1,
		}, {
			Name: "Group",
			URL:  "https://www.hentaicloud.com/videos/anal",
			Want: 23,
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
