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
			url:  "https://hentai.pro/ova-youkoso-sukebe-elf-no-mori-e-episode-4/",
			want: 1,
		}, {
			name: "Overview hentaistream.xxx",
			url:  "https://hentai.pro/tag/breasts/",
			want: 50,
		}, {
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
			name: "Single Episode hentai.pro",
			url:  "https://hentai.pro/bitch-na-inane-sama-episode-2/",
			want: 1,
		}, /* {
			name: "Overview hentai.pro",
			url:  "https://hentai.pro/tag/breasts/",
			want: 50,
		},*/{
			name: "Single Episode hentaistream.xxx",
			url:  "https://hentaistream.xxx/watch/netoshisu-episode-1_VpauMk84OoQXof3.html",
			want: 1,
		}, /* {
			name: "Overview hentaistream.xxx",
			url:  "https://hentaistream.xxx/videos/category/749",
			want: 47,
		},*/{
			name: "Single Episode hentaihaven.red/",
			url:  "https://hentaihaven.red/hentai/joshi-luck-episode-1/",
			want: 1,
		}, /* {
			name: "Overview hentaihaven.red/",
			url:  "https://hentaihaven.red/genre/2019-english/",
			want: 4,
			//can be more videos at the time when I am adding this it was only 4 -> normally it is 30 per site but that would be too much for testing
		}, {
			name: "Single Episode hentai.tv",
			url:  "https://hentai.tv/hentai/chiisana-tsubomi-no-sono-oku-ni-episode-1/",
			want: 1,
		}, {
			name: "Overview hentai.tv",
			url:  "https://hentai.tv/trending/",
			want: 24,
		},*/{
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
