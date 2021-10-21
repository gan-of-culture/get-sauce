package hentaihaven

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Episode",
			url:  "https://hentaihaven.xxx/watch/ero-konbini-tenchou/episode-4/",
			want: 1,
		}, {
			name: "Series",
			url:  "https://hentaihaven.xxx/watch/ero-konbini-tenchou/",
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
			name: "Single Episode",
			url:  "https://hentaihaven.xxx/watch/showtime-uta-no-onee-san-datte-shitai/episode-3/",
			want: 1,
		}, {
			name: "Series",
			url:  "https://hentaihaven.xxx/watch/ero-konbini-tenchou/",
			want: 4,
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
