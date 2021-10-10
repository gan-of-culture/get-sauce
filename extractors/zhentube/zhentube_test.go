package zhentube

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Episode",
			url:  "https://zhentube.com/soshite-watashi-wa-ojisan-ni-episode-4/",
			want: 1,
		}, {
			name: "Category",
			url:  "https://zhentube.com/category/2021/",
			want: 30,
		}, {
			name: "Tag",
			url:  "https://zhentube.com/tag/new-hentai-stream/",
			want: 30,
		}, {
			name: "Actor",
			url:  "https://zhentube.com/actor/kotomi/",
			want: 3,
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
			url:  "https://zhentube.com/torokase-orgasm-episode-1/",
			want: 1,
		}, /*{
			name: "Category",
			url:  "https://zhentube.com/category/censored-hentai/",
			want: 30,
		},*/
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
