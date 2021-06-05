package damn

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Episode",
			url:  "https://www.damn.stream/watch/hentai/reikan-shoujo-gaiden-toilet-no-hanako-san-vs-kukkyou-taimashi-akuochi-manko-ni-tenchuu-semen-renzoku-nakadashi-episode-1",
			want: 1,
		}, {
			name: "Series",
			url:  "https://www.damn.stream/hentai/reikan-shoujo-gaiden-toilet-no-hanako-san-vs-kukkyou-taimashi-akuochi-manko-ni-tenchuu-semen-renzoku-nakadashi",
			want: 2,
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
			url:  "https://www.damn.stream/watch/hentai/reikan-shoujo-gaiden-toilet-no-hanako-san-vs-kukkyou-taimashi-akuochi-manko-ni-tenchuu-semen-renzoku-nakadashi-episode-1",
			want: 1,
		}, {
			name: "Series",
			url:  "https://www.damn.stream/hentai/reikan-shoujo-gaiden-toilet-no-hanako-san-vs-kukkyou-taimashi-akuochi-manko-ni-tenchuu-semen-renzoku-nakadashi",
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
