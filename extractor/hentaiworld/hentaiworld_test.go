package hentaiworld

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Episode",
			url:  "https://hentaiworld.tv/hentai-videos/ijirare-fukushuu-saimin-episode-2/",
			want: 1,
		}, {
			name: "All episodes page",
			url:  "https://hentaiworld.tv/all-episodes/page/2/",
			want: 30,
		}, {
			name: "Uncensored page",
			url:  "https://hentaiworld.tv/uncensored/",
			want: 30,
		}, {
			name: "3d page",
			url:  "https://hentaiworld.tv/3d/",
			want: 60,
		}, {
			name: "tag page",
			url:  "https://hentaiworld.tv/hentai-videos/tag/anal/",
			want: 30,
		}, {
			name: "3d post page",
			url:  "https://hentaiworld.tv/hentai-videos/3d/final-fantasy-tifa-7/",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urls, err := ParseURL(tt.url)
			if err != nil {
				t.Error(err)
			}
			if len(urls) > tt.want {
				t.Errorf("Got: %v - want: %v", len(urls), tt.want)
			}
		})
	}
}

func TestExtractData(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "Single 3d extraction",
			url:  "https://hentaiworld.tv/hentai-videos/3d/final-fantasy-tifa-7/",
			want: "Final Fantasy – Tifa",
		},
		{
			name: "Single default extraction",
			url:  "https://hentaiworld.tv/hentai-videos/yuutousei-ayaka-no-uraomote-episode-1/",
			want: "Yuutousei Ayaka no Uraomote – Episode 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := ExtractData(tt.url)
			if data.Err != nil {
				t.Error(data.Err)
			}
			if data.Title != tt.want {
				t.Errorf("Got: %v - want: %v", data.Title, tt.want)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Page Extraction",
			url:  "https://hentaiworld.tv/hentai-videos/category/yuutousei-ayaka-no-uraomote/",
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Extract(tt.url)
			if err != nil {
				t.Error(err)
			}
			if len(data) > tt.want {
				t.Errorf("Got: %v - want: %v", len(data), tt.want)
			}
		})
	}
}
