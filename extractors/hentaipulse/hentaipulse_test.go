package hentaipulse

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Episode",
			url:  "https://hentaipulse.com/toshoshitsu-no-kanojo-seiso-na-kimi-ga-ochiru-made-the-animation-episode-04-english-subbed/",
			want: 1,
		}, {
			name: "Overview",
			url:  "https://hentaipulse.com/hentai-anime/english-subbed-hentai-anime/",
			want: 24,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urls := parseURL(tt.url)
			if len(urls) != tt.want || len(urls) == 0 {
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
			url:  "https://hentaipulse.com/toshoshitsu-no-kanojo-seiso-na-kimi-ga-ochiru-made-the-animation-episode-04-english-subbed/",
			want: 1,
		}, {
			name: "Overview",
			url:  "https://hentaipulse.com/hentai-anime/english-subbed-hentai-anime/",
			want: 24,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := New().Extract(tt.url)
			if err != nil {
				t.Error(err)
			}
			if len(data) != tt.want || len(data) == 0 {
				t.Errorf("Got: %v - want: %v", len(data), tt.want)
			}
		})
	}
}
