package simplyhentai

import (
	"net/url"
	"testing"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Gallery doujin.sexy",
			url:  "https://doujin.sexy/fate-grand-order/fdo-fatedosukebe-order-vol80",
			want: 1,
		}, {
			name: "Overview doujin.sexy",
			url:  "https://doujin.sexy/character/gudao",
			want: 24,
		}, {
			name: "Single Gallery simply-hentai.com",
			url:  "https://www.simply-hentai.com/1-kimetsu-no-yaiba/saimin-onsen-kanroji-mitsuri",
			want: 1,
		}, {
			name: "Overview simply-hentai.com",
			url:  "https://www.simply-hentai.com/tag/big-breasts",
			want: 24,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.url)
			if err != nil {
				t.Error(err)
			}

			site = "https://" + u.Host + "/"

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
			name: "Single Gallery simply-hentai.com",
			url:  "https://www.simply-hentai.com/original-work/torotoro-ni-shite-ageru-ch1-3",
			want: 1,
		}, {
			name: "Single Gallery doujin.sexy",
			url:  "https://doujin.sexy/fate-grand-order/fdo-fatedosukebe-order-vol80",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := New().Extract(tt.url)
			if err != nil {
				t.Error(err)
			}
			if len(data) > tt.want || len(data) == 0 {
				t.Errorf("Got: %v - want: %v", len(data), tt.want)
			}
		})
	}
}
