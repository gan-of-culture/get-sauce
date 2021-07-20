package muchohentai

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Episode",
			url:  "https://muchohentai.com/aBo4Rk/158393/",
			want: 1,
		}, {
			name: "Genre",
			url:  "https://muchohentai.com/g/1080p/",
			want: 24,
		}, {
			name: "Series",
			url:  "https://muchohentai.com/s/overflow/",
			want: 24,
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
			url:  "https://muchohentai.com/aBo4Rk/158393/",
			want: 24,
		}, {
			name: "Single Episode Single Stream",
			url:  "https://muchohentai.com/aBo4Rk/161377/",
			want: 1,
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
