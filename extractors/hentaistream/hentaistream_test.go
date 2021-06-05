package hentaistream

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Episode",
			url:  "https://hentaistream.moe/593/kateikyoushi-no-oneesan-2-the-animation-1/",
			want: 1,
		}, {
			name: "Series",
			url:  "https://hentaistream.moe/anime/kateikyoushi-no-oneesan-2-the-animation/",
			want: 2,
		}, {
			name: "Single Episode 4k",
			url:  "https://hentaistream.moe/515/overflow-1/",
			want: 1,
		}, {
			name: "Series 4k",
			url:  "https://hentaistream.moe/anime/overflow/",
			want: 8,
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
			name: "Series Extraction 4k",
			url:  "https://hentaistream.moe/anime/overflow/",
			want: 8,
		}, {
			name: "Single Episode 4k",
			url:  "https://hentaistream.moe/515/overflow-1/",
			want: 1,
		}, {
			name: "Single Episode",
			url:  "https://hentaistream.moe/593/kateikyoushi-no-oneesan-2-the-animation-1/",
			want: 1,
		}, {
			name: "Series",
			url:  "https://hentaistream.moe/anime/kateikyoushi-no-oneesan-2-the-animation/",
			want: 2,
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
