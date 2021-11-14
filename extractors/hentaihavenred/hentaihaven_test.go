package hentaihavenred

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Episode hentaihaven.red/",
			url:  "https://hentaihaven.red/hentai/joshi-luck-episode-1/",
			want: 1,
		}, {
			name: "Overview hentaihaven.red/",
			url:  "https://hentaihaven.red/ratings/",
			want: 35,
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
			name: "Single Episode hentaihaven.red",
			url:  "https://hentaihaven.red/hentai/bitch-na-inane-sama-episode-4/",
			want: 1,
		}, {
			name: "[OLD] Single Episode hentaihaven.red",
			url:  "https://hentaihaven.red/hentai/mako-chan-kaihatsu-nikki-episode-2/",
			want: 1,
		}, {
			name: "Overview hentaihaven.red",
			url:  "https://hentaihaven.red/genre/2019-english/",
			want: 4,
			//can be more videos at the time when I am adding this it was only 4 -> normally it is 30 per site but that would be too much for testing
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
