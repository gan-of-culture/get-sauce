package booru

import (
	"log"
	"testing"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want string
	}{
		{
			Name: "Tag query",
			URL:  "https://booru.io/q/1girl%20nude%20animal_ears%20cat%20solo",
			Want: "https://booru.io/api/legacy/query/entity?query=1girl%20nude%20animal_ears%20cat%20solo",
		}, {
			Name: "Single Tag query",
			URL:  "https://booru.io/q/1girl",
			Want: "https://booru.io/api/legacy/query/entity?query=1girl",
		}, {
			Name: "Example Post",
			URL:  "https://booru.io/p/YoZR3jurfVNOXD4vjCNn",
			Want: "https://booru.io/api/legacy/entity/YoZR3jurfVNOXD4vjCNn",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			log.Println(tt.Name)
			URL, err := parseURL(tt.URL)
			if err != nil {
				t.Error(err)
			}
			if URL != tt.Want {
				t.Errorf("Got: %v - Want: %v", URL, tt.Want)
			}
		})
	}
}

func TestExtractData(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Default extraction",
			URL:  "https://booru.io/api/legacy/entity/YoZR3jurfVNOXD4vjCNn",
			Want: 1,
		},
		{
			Name: "Query extraction",
			URL:  "https://booru.io/api/legacy/query/entity?query=1girl%20nude%20animal_ears%20cat%20solo",
			Want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := extractData(tt.URL)
			if err != nil {
				t.Error(err)
			}
			if len(data) > tt.Want {
				t.Errorf("Got: %v - Want: %v", len(data), tt.Want)
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
			Name: "Default extraction",
			URL:  "https://booru.io/p/YoZR3jurfVNOXD4vjCNn",
			Want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.URL)
			if err != nil {
				t.Error(err)
			}
			if len(data) < tt.Want {
				t.Errorf("Got: %v - Want: %v", len(data), tt.Want)
			}
		})
	}
}
