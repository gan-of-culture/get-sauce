package booruproject

import (
	"testing"

	"github.com/gan-of-culture/go-hentai-scraper/config"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		amount int
		want   int
	}{
		{
			name: "Single image",
			url:  "https://rule34.xxx/index.php?page=post&s=view&id=4470590",
			want: 1,
		}, {
			name: "Overview page",
			url:  "https://rule34.xxx/index.php?page=post&s=list&tags=all",
			// atleast more than 2
			want: 2,
		}, {
			name: "Mass extract page",
			url:  "https://rule34.xxx/index.php?page=post&s=list&tags=all",
			// atleast more than 2
			want: 100,
		}, {
			name: "Mass extract page",
			url:  "https://gelbooru.com/index.php?page=post&s=list&tags=all",
			// atleast more than 2
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Mass extract page" {
				config.Amount = 101
			}
			elements := ParseURL(tt.url)

			if len(elements) < tt.want {
				t.Errorf("Got: %v - want: %v", len(elements), tt.want)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	type want struct {
		Title   string
		Type    string
		DataLen int
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Test image",
			url:  "https://gelbooru.com/index.php?page=post&s=view&id=5888807",
			want: want{
				Title:   "gelbooru_5888807",
				Type:    "image/png",
				DataLen: 1,
			},
		}, {
			name: "Test video",
			url:  "https://rule34.xxx/index.php?page=post&s=view&id=4470579",
			want: want{
				Title:   "rule34_4470579",
				Type:    "video/mp4",
				DataLen: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Amount = 26
			elements, err := Extract(tt.url)
			if err != nil {
				t.Error("elements has error or is too big for single tests")
			}
			act := want{
				Title:   elements[0].Title,
				Type:    elements[0].Type,
				DataLen: len(elements),
			}
			if act != tt.want {
				t.Errorf("Got: %v - want: %v", act, tt.want)
			}
		})
	}
}

func TestMassExtract(t *testing.T) {

	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Test mass",
			url:  "https://rule34.xxx/index.php?page=post&s=list&tags=all",
			want: 62,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Amount = 62
			elements, err := Extract(tt.url)
			if err != nil {
				t.Error("elements has error or is too big for single tests")
			}
			if len(elements) != tt.want {
				t.Errorf("Got: %v - want: %v", len(elements), tt.want)
			}
		})
	}
}
