package htdoujin

import (
	"fmt"
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
			name: "Single Gallery HentaiEra",
			url:  "https://hentaiera.com/gallery/150354/",
			want: 1,
		}, {
			name: "Tag HentaiEra",
			url:  "https://hentaiera.com/tag/ahegao/",
			want: 25,
		}, {
			name: "Single Gallery IMHentai",
			url:  "https://imhentai.xxx/gallery/684976/",
			want: 1,
		}, {
			name: "Tag IMHentai",
			url:  "https://imhentai.xxx/artist/100yen-locker/",
			want: 20,
		}, {
			name: "Single Gallery HentaiFox",
			url:  "https://hentaifox.com/gallery/84580/",
			want: 1,
		}, {
			name: "Tag HentaiFox",
			url:  "https://hentaifox.com/tag/age-progression/",
			want: 20,
		}, {
			name: "Single Gallery HentaiEra",
			url:  "https://hentaiera.com/gallery/150354/",
			want: 1,
		}, {
			name: "Tag HentaiEra",
			url:  "https://hentaiera.com/tag/ahegao/",
			want: 25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.url)
			if err != nil {
				t.Error(err)
			}

			if cdnPrefix, ok := sites[u.Host]; ok {
				site = "https://" + u.Host + "/"
				cdn = fmt.Sprintf("https://%s.%s/", cdnPrefix, u.Host)
			}

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
			name: "Single Gallery IMHentai",
			url:  "https://imhentai.xxx/gallery/684976/",
			want: 1,
		}, {
			name: "Single Gallery HentaiFox",
			url:  "https://hentaifox.com/gallery/84580/",
			want: 1,
		}, {
			name: "Single Gallery HentaiEra",
			url:  "https://hentaiera.com/gallery/488946/",
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
