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
			name: "Single Episode uncensoredhentai.xxx",
			url:  "https://uncensoredhentai.xxx/watch/mako-chan-kaihatsu-nikki-episode-1/",
			want: 1,
		}, {
			name: "Overview uncensoredhentai.xxx",
			url:  "https://uncensoredhentai.xxx/genres/ahegao/",
			want: 18,
		}, {
			name: "Single Episode hentai.pro",
			url:  "https://hentai.pro/knight-of-erin-episode-2/",
			want: 1,
		}, {
			name: "Overview hentai.pro",
			url:  "https://hentai.pro/tag/breasts/",
			want: 50,
		}, {
			name: "Single Episode hentaistream.xxx",
			url:  "https://hentaistream.xxx/watch/ijirare-fukushuu-saimin-episode-1/",
			want: 1,
		}, {
			name: "Overview hentaistream.xxx",
			url:  "https://hentaistream.xxx/genres/ahegao/",
			want: 18,
		}, /*{
			name: "Single Episode hentai.tv",
			url:  "https://hentai.tv/hentai/chiisana-tsubomi-no-sono-oku-ni-episode-1/",
			want: 1,
		}, {
			name: "Overview hentai.tv",
			url:  "https://hentai.tv/trending/",
			want: 24,
		}, */{
			name: "Single Episode animeidhentai.com",
			url:  "https://animeidhentai.com/31678/mako-chan-kaihatsu-nikki-episode-1/",
			want: 1,
		}, {
			name: "Series animeidhentai.com",
			url:  "https://animeidhentai.com/hentai/mako-chan-kaihatsu-nikki/",
			want: 4,
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
			name: "Single Episode uncensoredhentai.xxx",
			url:  "https://uncensoredhentai.xxx/watch/mako-chan-kaihatsu-nikki-episode-1/",
			want: 1,
		}, {
			name: "Single Episode hentai.pro",
			url:  "https://hentai.pro/imaizumin-chi-wa-douyara-gal-no-tamariba-ni-natteru-rashii-episode-2/",
			want: 1,
		}, {
			name: "Single Episode hentaistream.xxx",
			url:  "https://hentaistream.xxx/watch/mako-chan-kaihatsu-nikki-episode-1/",
			want: 1,
		}, /* {
			name: "Single Episode hentai.tv",
			url:  "https://hentai.tv/hentai/chiisana-tsubomi-no-sono-oku-ni-episode-1/",
			want: 1,
		}, */{
			name: "Single Episode animeidhentai.com",
			url:  "https://animeidhentai.com/31680/mako-chan-kaihatsu-nikki-episode-2/",
			want: 1,
		}, {
			name: "Series animeidhentai.com",
			url:  "https://animeidhentai.com/hentai/mako-chan-kaihatsu-nikki/",
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := New().Extract(tt.url)
			if err != nil && !strings.Contains(err.Error(), "Video not found") {
				t.Error(err)
			}
			if len(data) > tt.want || len(data) == 0 {
				t.Errorf("Got: %v - want: %v", len(data), tt.want)
			}
		})
	}
}
