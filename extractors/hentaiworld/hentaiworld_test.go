package hentaiworld

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single Episode",
			URL:  "https://hentaiworld.tv/hentai-videos/ijirare-fukushuu-saimin-episode-2/",
			Want: 1,
		}, {
			Name: "All episodes page",
			URL:  "https://hentaiworld.tv/all-episodes/page/2/",
			Want: 30,
		}, {
			Name: "Uncensored page",
			URL:  "https://hentaiworld.tv/uncensored/",
			Want: 30,
		}, {
			Name: "3d page",
			URL:  "https://hentaiworld.tv/3d/",
			Want: 60,
		}, {
			Name: "tag page",
			URL:  "https://hentaiworld.tv/hentai-videos/tag/anal/",
			Want: 30,
		}, {
			Name: "3d post page",
			URL:  "https://hentaiworld.tv/hentai-videos/3d/final-fantasy-tifa-7/",
			Want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			URLs := parseURL(tt.URL)
			if len(URLs) > tt.Want {
				t.Errorf("Got: %v - Want: %v", len(URLs), tt.Want)
			}
		})
	}
}

func TestExtractData(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want string
	}{
		{
			Name: "Single 3d extraction",
			URL:  "https://hentaiworld.tv/hentai-videos/3d/final-fantasy-tifa-7/",
			Want: "Final Fantasy – Tifa",
		},
		{
			Name: "Single default extraction",
			URL:  "https://hentaiworld.tv/hentai-videos/yuutousei-ayaka-no-uraomote-episode-1/",
			Want: "Yuutousei Ayaka no Uraomote - Episode 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := extractData(tt.URL)
			if err != nil {
				t.Error(err)
			}
			if data.Title != tt.Want {
				t.Errorf("Got: %v - Want: %v", data.Title, tt.Want)
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
			Name: "Page Extraction",
			URL:  "https://hentaiworld.tv/hentai-videos/category/yuutousei-ayaka-no-uraomote/",
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
