package iwara

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/config"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single video",
			url:  "https://ecchi.iwara.tv/videos/kmnzvsa75uzbaw36?language=en",
			want: 1,
		}, {
			name: "Single images",
			url:  "https://ecchi.iwara.tv/images/%E6%B9%AF%E4%B8%8A%E3%81%8C%E3%82%8A%E3%82%86%E3%81%84%E3%81%A1%E3%82%83%E3%82%93?language=en",
			want: 2,
		}, /*{
			name: "Mass",
			url:  "https://ecchi.iwara.tv/images?language=en&f%5B0%5D=field_image_categories%3A5&page=1",
			want: 40,
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Mass" {
				config.Amount = 40
			}

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
			name: "Single video",
			url:  "https://ecchi.iwara.tv/videos/kmnzvsa75uzbaw36?language=en",
			want: 1,
		}, {
			name: "Single images",
			url:  "https://ecchi.iwara.tv/images/%E6%B9%AF%E4%B8%8A%E3%81%8C%E3%82%8A%E3%82%86%E3%81%84%E3%81%A1%E3%82%83%E3%82%93?language=en",
			want: 2,
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
