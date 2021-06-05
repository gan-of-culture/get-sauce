package hentaidude

import (
	"testing"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Episode",
			url:  "https://hentaidude.com/aisei-tenshi-love-mary-akusei-jutai-episode-1/",
			want: 1,
		}, {
			name: "Landing page",
			url:  "https://hentaidude.com/",
			want: 20,
		}, {
			name: "Orderby",
			url:  "https://hentaidude.com/?orderby=date",
			want: 20,
		}, {
			name: "Tags",
			url:  "https://hentaidude.com/?orderby=date&tid=1472",
			want: 20,
		}, {
			name: "Different page",
			url:  "https://hentaidude.com/page/3/?orderby=date&tid=1472",
			want: 20,
		}, {
			name: "3D Single Episode",
			url:  "https://hentaidude.com/scarlet-nights-episode-1/",
			want: 1,
		}, {
			name: "3D Landing page",
			url:  "https://hentaidude.com/tag/3d-hentai-0/",
			want: 20,
		}, {
			name: "3D Orderby",
			url:  "https://hentaidude.com/tag/3d-hentai-0/?orderby=date",
			want: 20,
		}, {
			name: "3D Tags",
			url:  "https://hentaidude.com/tag/3d-hentai-0/?orderby=date&tid=1472",
			want: 20,
		}, {
			name: "3D Different page",
			url:  "https://hentaidude.com/tag/3d-hentai-0/page/2/?tid=1541",
			want: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urls := parseURL(tt.url)
			if len(urls) > tt.want {
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
			name: "Page",
			url:  "https://hentaidude.com/",
			want: 20,
		}, {
			name: "Single Episode",
			url:  "https://hentaidude.com/yuutousei-ayaka-no-uraomote-episode-1/",
			want: 1,
		}, {
			name: "3D Single Episode",
			url:  "https://hentaidude.com/scarlet-nights-episode-1/",
			want: 1,
		}, /*{
			name: "3D Page",
			url:  "https://hentaidude.com/tag/3d-hentai-0/",
			want: 20,
		},*/ //this was removed to save time when testing
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
