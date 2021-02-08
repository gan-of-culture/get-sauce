package hentais

import (
	"testing"
)

func TestExtractData(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single default extraction",
			url:  "https://www.hentais.tube/episodes/shishunki-sex-episode-4/",
			want: 2,
		},
		{
			name: "Whole default series extraction",
			url:  "https://www.hentais.tube/tvshows/shishunki-sex/",
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := ExtractData(tt.url)
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
			name: "Single default extraction",
			url:  "https://www.hentais.tube/episodes/shishunki-sex-episode-4",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Extract(tt.url)
			if err != nil {
				t.Error(err)
			}
			if len(data) < tt.want {
				t.Errorf("Got: %v - want: %v", len(data), tt.want)
			}
		})
	}
}
