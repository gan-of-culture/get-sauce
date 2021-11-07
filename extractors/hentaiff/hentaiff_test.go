package hentaiff

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Episode Raw",
			url:  "https://hentaiff.com/onaho-kyoushitsu-joshi-zenin-ninshin-keikaku-the-animation-raw/",
			want: 1,
		}, {
			name: "Single Episode Eng Sub",
			url:  "https://hentaiff.com/onaho-kyoushitsu-joshi-zenin-ninshin-keikaku-the-animation-english-subbed/",
			want: 1,
		}, {
			name: "Single Episode Eng Dub",
			url:  "https://hentaiff.com/a-kite-episode-02-english-dubbed/",
			want: 1,
		}, {
			name: "Single Episode Preview",
			url:  "https://hentaiff.com/onaho-kyoushitsu-joshi-zenin-ninshin-keikaku-the-animation-previews/",
			want: 1,
		}, {
			name: "Series",
			url:  "https://hentaiff.com/anime/a-kite/",
			want: 6,
		}, {
			// this is the same logic for all extensions that group shows e.g. /genres/
			// its hard to make a test for the other groups since the number of episodes always changes
			name: "Studio",
			url:  "https://hentaiff.com/studio/arms/",
			// 5 show with a sum of 28 episodes
			want: 28,
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
			url:  "https://hentaiff.com/onaho-kyoushitsu-joshi-zenin-ninshin-keikaku-the-animation-english-subbed/",
			want: 1,
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
