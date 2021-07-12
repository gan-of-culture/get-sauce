package booru

import (
	"log"
	"testing"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "Tag query",
			url:  "https://booru.io/q/1girl%20nude%20animal_ears%20cat%20solo",
			want: "https://booru.io/api/query/entity?query=1girl%20nude%20animal_ears%20cat%20solo",
		}, {
			name: "Example Post",
			url:  "https://booru.io/p/YoZR3jurfVNOXD4vjCNn",
			want: "https://booru.io/api/entity/YoZR3jurfVNOXD4vjCNn",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Println(tt.name)
			url, err := parseURL(tt.url)
			if err != nil {
				t.Error(err)
			}
			if url != tt.want {
				t.Errorf("Got: %v - want: %v", url, tt.want)
			}
		})
	}
}

func TestExtractData(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Default extraction",
			url:  "https://booru.io/api/entity/YoZR3jurfVNOXD4vjCNn",
			want: 1,
		},
		{
			name: "Query extraction",
			url:  "https://booru.io/api/query/entity?query=1girl%20nude%20animal_ears%20cat%20solo",
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := extractData(tt.url)
			if err != nil {
				t.Error(err)
			}
			if len(data) > tt.want {
				t.Errorf("Got: %v - want: %v", len(data), tt.want)
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
			name: "Default extraction",
			url:  "https://booru.io/p/YoZR3jurfVNOXD4vjCNn",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := New().Extract(tt.url)
			if err != nil {
				t.Error(err)
			}
			if len(data) < tt.want {
				t.Errorf("Got: %v - want: %v", len(data), tt.want)
			}
		})
	}
}
