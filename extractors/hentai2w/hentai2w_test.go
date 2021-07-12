package hentai2w

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Video",
			url:  "https://hentai2w.com/video/youkoso-sukebe-elf-no-mori-e-episode-2-3693.html",
			want: 1,
		}, {
			name: "Category",
			url:  "https://hentai2w.com/channels/125/magic/",
			want: 40,
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
			name: "Single Video",
			url:  "https://hentai2w.com/video/youkoso-sukebe-elf-no-mori-e-episode-2-3693.html",
			want: 1,
		}, /*{
			name: "Category",
			url:  "https://hentai2w.com/channels/125/magic/",
			want: 40,
		},*/
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
