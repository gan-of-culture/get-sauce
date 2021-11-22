package iwara

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/config"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single video",
			URL:  "https://ecchi.iwara.tv/videos/kmnzvsa75uzbaw36?language=en",
			Want: 1,
		}, {
			Name: "Single images",
			URL:  "https://ecchi.iwara.tv/images/%E6%B9%AF%E4%B8%8A%E3%81%8C%E3%82%8A%E3%82%86%E3%81%84%E3%81%A1%E3%82%83%E3%82%93?language=en",
			Want: 2,
		}, /*{
			Name: "Mass",
			URL:  "https://ecchi.iwara.tv/images?language=en&f%5B0%5D=field_image_categories%3A5&page=1",
			Want: 40,
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.Name == "Mass" {
				config.Amount = 40
			}

			URLs := parseURL(tt.URL)
			if len(URLs) > tt.Want || len(URLs) == 0 {
				t.Errorf("Got: %v - Want: %v", len(URLs), tt.Want)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single video",
			URL:  "https://ecchi.iwara.tv/videos/kmnzvsa75uzbaw36?language=en",
			Want: 1,
		}, {
			Name: "Single images",
			URL:  "https://ecchi.iwara.tv/images/%E6%B9%AF%E4%B8%8A%E3%81%8C%E3%82%8A%E3%82%86%E3%81%84%E3%81%A1%E3%82%83%E3%82%93?language=en",
			Want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.URL)
			if err != nil {
				t.Error(err)
			}
			if len(data) > tt.Want || len(data) == 0 {
				t.Errorf("Got: %v - Want: %v", len(data), tt.Want)
			}
		})
	}
}
