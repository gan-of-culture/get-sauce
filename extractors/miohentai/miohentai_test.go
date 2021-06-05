package miohentai

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Episode",
			url:  "https://miohentai.com/video/enjo-kouhai-episode-2/",
			want: 1,
		}, {
			name: "Tag",
			url:  "https://miohentai.com/tag/1080p/",
			want: 22,
		}, {
			name: "Image",
			url:  "https://miohentai.com/image-library/the-latest-influencers-2020-dress/",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urls := parseURL(tt.url)
			if len(urls) > tt.want || len(urls) == 0 {
				t.Errorf("Got: %v - want: %v", len(urls), tt.want)
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
			name: "Single Episode",
			url:  "https://miohentai.com/video/enjo-kouhai-episode-2/",
			want: 1,
		}, {
			name: "Tag",
			url:  "https://miohentai.com/tag/1080p/",
			want: 22,
		}, {
			name: "Image",
			url:  "https://miohentai.com/image-library/the-latest-influencers-2020-dress/",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := New().Extract(tt.url)
			if err != nil {
				t.Error(err)
			}
			if len(data) > tt.want {
				t.Errorf("Got: %v - want: %v", len(data), tt.want)
			}
		})
	}
}
