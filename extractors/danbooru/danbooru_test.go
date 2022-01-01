package danbooru

import (
	"log"
	"testing"

	"github.com/gan-of-culture/get-sauce/test"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Overview page",
			URL:  "https://danbooru.donmai.us/posts?page=3&tags=fire_emblem",
			Want: 2,
		}, {
			Name: "Example Post",
			URL:  "https://danbooru.donmai.us/posts/3749687",
			Want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			log.Println(tt.Name)
			URLs, err := parseURL(tt.URL)
			if err != nil {
				t.Error(err)
			}
			if len(URLs) < tt.Want {
				t.Errorf("Got: %v - Want: %v", len(URLs), tt.Want)
			}
		})
	}
}

func TestExtractData(t *testing.T) {
	type Want struct {
		numberOfStream int
		title          string
	}
	tests := []struct {
		Name string
		URL  string
		Want Want
	}{
		{
			Name: "Default extraction",
			URL:  "https://danbooru.donmai.us/posts/3773519",
			Want: Want{
				numberOfStream: 1,
				title:          "misty and squirtle (pokemon and 2 more) drawn by shellvi",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := extractData(tt.URL)
			if err != nil {
				t.Error(err)
			}
			if len(data.Streams) != tt.Want.numberOfStream {
				t.Errorf("Got: %v - Want: %v", len(data.Streams), tt.Want.numberOfStream)
			}
			if data.Title != tt.Want.title {
				t.Errorf("Got: %v - Want: %v", data.Title, tt.Want.title)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		Name string
		Args test.Args
	}{
		{
			Name: "Default extraction",
			Args: test.Args{
				URL:     "https://danbooru.donmai.us/posts/3749687",
				Title:   "konpaku youmu and konpaku youmu (touhou) drawn by niwashi_(yuyu)",
				Quality: "1782 x 2048",
				Size:    157584,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.Args.URL)
			test.CheckError(t, err)
			test.Check(t, tt.Args, data[0])
		})
	}
}
