package hanime

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Episode",
			url:  "https://hanime.tv/videos/hentai/toilet-no-hanako-san-vs-kukkyou-taimashi-2",
			want: 1,
		}, {
			name: "Category",
			url:  "https://hanime.tv/browse/tags/ahegao",
			want: 24,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urls := ParseURL(tt.url)
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
			url:  "https://hanime.tv/videos/hentai/toilet-no-hanako-san-vs-kukkyou-taimashi-2",
			want: 1,
		}, {
			name: "Series",
			url:  "https://hanime.tv/browse/tags/ahegao",
			want: 24,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Extract(tt.url)
			if err != nil {
				t.Error(err)
			}
			if len(data) > tt.want {
				t.Errorf("Got: %v - want: %v", len(data), tt.want)
			}
		})
	}
}
