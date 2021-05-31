package hentai2read

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Gallery",
			url:  "https://hentai2read.com/okitasan_to_kotasu_ecchi/#availableChapters",
			want: 1,
		}, {
			name: "Tag",
			url:  "https://hentai2read.com/hentai-list/category/Romance/",
			want: 48,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urls := ParseURL(tt.url)
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
			name: "Single Gallery",
			url:  "https://hentai2read.com/elevenpm_miniature_garden/#availableChapters",
			want: 1,
		}, /*{
			name: "Tag",
			url:  "https://hentai2read.com/hentai-list/category/Romance/",
			want: 20,
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Extract(tt.url)
			if err != nil {
				t.Error(err)
			}
			if len(data) > tt.want || len(data) == 0 {
				t.Errorf("Got: %v - want: %v", len(data), tt.want)
			}
		})
	}
}
