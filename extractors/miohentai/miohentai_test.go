package miohentai

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single Episode",
			URL:  "https://miohentai.com/video/enjo-kouhai-episode-2/",
			Want: 1,
		}, {
			Name: "Tag",
			URL:  "https://miohentai.com/tag/1080p/",
			Want: 22,
		}, {
			Name: "Image",
			URL:  "https://miohentai.com/image-library/the-latest-influencers-2020-dress/",
			Want: 1,
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
			URL:  "https://miohentai.com/video/enjo-kouhai-episode-2/",
			Want: 1,
		}, {
			Name: "Tag",
			URL:  "https://miohentai.com/tag/1080p/",
			Want: 22,
		}, {
			Name: "Image",
			URL:  "https://miohentai.com/image-library/the-latest-influencers-2020-dress/",
			Want: 1,
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
