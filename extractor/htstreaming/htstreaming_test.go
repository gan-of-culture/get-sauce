package htstreaming

import (
	"strings"
	"testing"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Episode hentaistream.xxx",
			url:  "https://hentaistream.xxx/watch/tonari-no-ie-no-anette-san-the-animation-episode-1_waYqxLSASjFPICZ.html",
			want: 1,
		}, {
			name: "Overview hentaistream.xxx",
			url:  "https://hentaistream.xxx/videos/category/749",
			want: 20,
		}, {
			name: "Single Episode hentaihaven.red/",
			url:  "https://hentaihaven.red/hentai/joshi-luck-episode-1/",
			want: 1,
		}, {
			name: "Overview hentaihaven.red/",
			url:  "https://hentaihaven.red/ratings/",
			want: 35,
		}, /*{
			name: "Single Episode hentai.tv",
			url:  "https://hentai.tv/hentai/chiisana-tsubomi-no-sono-oku-ni-episode-1/",
			want: 1,
		}, {
			name: "Overview hentai.tv",
			url:  "https://hentai.tv/trending/",
			want: 24,
		},*/{
			name: "Single Episode animeidhentai.com",
			url:  "https://animeidhentai.com/31577/doukyuusei-natsu-no-owari-ni-episode-1/",
			want: 1,
		}, {
			name: "Series animeidhentai.com",
			url:  "https://animeidhentai.com/hentai/doukyuusei-natsu-no-owari-ni/",
			want: 4,
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
			name: "Single Episode hentaistream.xxx",
			url:  "https://hentaistream.xxx/watch/ecchi-na-onee-chan-ni-shiboraretai-episode-1-subbed_2v5zblbKJSynGJ6.html",
			want: 1,
		}, {
			name: "Overview hentaistream.xxx",
			url:  "https://hentaistream.xxx/videos/category/749",
			want: 20,
		}, {
			name: "Single Episode hentaihaven.red/",
			url:  "https://hentaihaven.red/hentai/joshi-luck-episode-1/",
			want: 1,
		}, {
			name: "Overview hentaihaven.red/",
			url:  "https://hentaihaven.red/ratings/",
			want: 35,
		}, /*{
			name: "Single Episode hentai.tv",
			url:  "https://hentai.tv/hentai/chiisana-tsubomi-no-sono-oku-ni-episode-1/",
			want: 1,
		}, {
			name: "Overview hentai.tv",
			url:  "https://hentai.tv/trending/",
			want: 24,
		},*/{
			name: "Single Episode animeidhentai.com",
			url:  "https://animeidhentai.com/31577/doukyuusei-natsu-no-owari-ni-episode-1/",
			want: 1,
		}, {
			name: "Series animeidhentai.com",
			url:  "https://animeidhentai.com/hentai/doukyuusei-natsu-no-owari-ni/",
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Extract(tt.url)
			if err != nil && !strings.Contains(err.Error(), "Video not found") {
				t.Error(err)
			}
			if len(data) > tt.want || len(data) == 0 {
				t.Errorf("Got: %v - want: %v", len(data), tt.want)
			}
		})
	}
}
