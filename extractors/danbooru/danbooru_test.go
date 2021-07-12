package danbooru

import (
	"log"
	"testing"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Overview page",
			url:  "https://danbooru.donmai.us/posts?page=3&tags=fire_emblem",
			want: 2,
		}, {
			name: "Example Post",
			url:  "https://danbooru.donmai.us/posts/3749687",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log.Println(tt.name)
			urls, err := parseURL(tt.url)
			if err != nil {
				t.Error(err)
			}
			if len(urls) < tt.want {
				t.Errorf("Got: %v - want: %v", len(urls), tt.want)
			}
		})
	}
}

func TestExtractData(t *testing.T) {
	type want struct {
		numberOfStream int
		title          string
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Default extraction",
			url:  "https://danbooru.donmai.us/posts/3773519",
			want: want{
				numberOfStream: 1,
				title:          "misty and squirtle (pokemon and 2 more) drawn by shellvi",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := extractData(tt.url)
			if err != nil {
				t.Error(err)
			}
			if len(data.Streams) != tt.want.numberOfStream {
				t.Errorf("Got: %v - want: %v", len(data.Streams), tt.want.numberOfStream)
			}
			if data.Title != tt.want.title {
				t.Errorf("Got: %v - want: %v", data.Title, tt.want.title)
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
			url:  "https://danbooru.donmai.us/posts/3749687",
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
