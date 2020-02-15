package exhentai

import (
	"testing"

	"github.com/gan-of-culture/go-hentai-scraper/config"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		URL  string
		//expect atlest ... galleries
		numberOfGalleries int
	}{
		{
			name:              "Single gallery",
			URL:               "https://exhentai.org/g/1566926/6c5691abf3/",
			numberOfGalleries: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			URLs := ParseURL(tt.URL)
			if len(URLs) < tt.numberOfGalleries {
				t.Errorf("Got: %v - want: %v", len(URLs), tt.numberOfGalleries)
			}
		})
	}
}
func TestExtract(t *testing.T) {
	tests := []struct {
		name string
		URL  string
		// expect atleast ... data structs
		numberOfData int
	}{
		{
			name:         "Single gallery",
			URL:          "https://exhentai.org/g/1566926/6c5691abf3/",
			numberOfData: 24,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Username = "idinonothuman"
			config.UserPassword = "idinonothuman08122019"
			data, err := Extract(tt.URL)
			if err != nil {
				t.Error(err)
			}
			if len(data) < tt.numberOfData {
				t.Errorf("Got: %v - want: %v", len(data), tt.numberOfData)
			}
		})
	}
}
