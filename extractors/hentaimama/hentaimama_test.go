package hentaimama

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Episode",
			url:  "https://hentaimama.io/episodes/katainaka-ni-totsui-de-kita-russia-musume-h-shimakuru-ohanashi-episode-1/",
			want: 1,
		}, {
			name: "Series",
			url:  "https://hentaimama.io/tvshows/katainaka-ni-totsui-de-kita-russia-musume-h-shimakuru-ohanashi/",
			want: 4,
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
			name: "Single Episode",
			url:  "https://hentaimama.io/episodes/ura-jutaijima-episode-1/",
			want: 1,
		}, {
			name: "Series",
			url:  "https://hentaimama.io/tvshows/ura-jutaijima/",
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
