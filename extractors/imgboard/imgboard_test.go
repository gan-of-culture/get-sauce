package imgboard

import (
	"strings"
	"testing"

	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/static"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Page",
			url:  "https://rule34.xxx/index.php?page=post&s=list&tags=world_of_warcraft+video+draenei&pid=378",
			// atleast more than 2
			want: 2,
		}, {
			name: "Mass extract page booru project",
			url:  "https://tbib.org/index.php?page=post&s=list&tags=1girl+solo+uncensored+full_body+&pid=0",
			want: 100,
		},
		{
			name: "Single image booru project",
			url:  "https://rule34.xxx/index.php?page=post&s=view&id=4470590",
			want: 1,
		}, {
			name: "Overview page booru project",
			url:  "https://rule34.xxx/index.php?page=post&s=list&tags=all",
			// atleast more than 2
			want: 2,
		}, {
			name: "Mass extract page booru project",
			url:  "https://rule34.xxx/index.php?page=post&s=list&tags=all",
			// atleast more than 2
			want: 100,
		}, {
			name: "Mass extract page booru project2",
			url:  "https://gelbooru.com/index.php?page=post&s=list&tags=all",
			// atleast more than 2
			want: 100,
		}, {
			name: "Single image",
			url:  "https://yande.re/post/show/745150",
			want: 1,
		}, {
			name: "Overview page booru project",
			url:  "https://konachan.com/post?tags=uncensored",
			// atleast more than 2
			want: 2,
		}, {
			name: "Mass extract page",
			url:  "https://konachan.com/post?tags=uncensored",
			// atleast more than 2
			want: 100,
		}, {
			name: "Mass extract page2",
			url:  "https://yande.re/post?tags=tateha&page=2",
			// atleast more than 2
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if strings.HasPrefix(tt.name, "Mass extract") {
				config.Amount = 101
			}
			elements := parseURL(tt.url)

			if len(elements) < tt.want {
				t.Errorf("Got: %v - want: %v", len(elements), tt.want)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	type want struct {
		Title   string
		Type    static.DataType
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
				Type:    static.DataTypeImage,
				DataLen: 1,
			},
		}, {
			name: "Test video",
			url:  "https://rule34.xxx/index.php?page=post&s=view&id=4470579",
			want: want{
				Title:   "rule34_4470579",
				Type:    static.DataTypeVideo,
				DataLen: 1,
			},
		}, {
			name: "Test image konachan",
			url:  "https://konachan.com/post/show/323560/anthropomorphism-azur_lane-black_hair-blush-breast",
			want: want{
				Title:   "konachan_323560",
				Type:    static.DataTypeImage,
				DataLen: 1,
			},
		}, {
			name: "Test image yande.re",
			url:  "https://yande.re/post/show/745150",
			want: want{
				Title:   "yande_745150",
				Type:    static.DataTypeImage,
				DataLen: 1,
			},
		}, {
			name: "Test image tbib",
			url:  "https://tbib.org/index.php?page=post&s=view&id=9022091",
			want: want{
				Title:   "tbib_9022091",
				Type:    static.DataTypeImage,
				DataLen: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elements, err := New().Extract(tt.url)
			if err != nil {
				t.Error(err)
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

func TestExtractDataFromDirectLink(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "DL extract yande.re",
			url:  "https://yande.re/post?",
			want: 10,
		}, {
			name: "DL extract konachan",
			url:  "https://konachan.com/post?tags=",
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Amount = 10
			elements, err := New().Extract(tt.url)
			if err != nil {
				t.Error(err)
			}
			if len(elements) < tt.want {
				t.Errorf("Got: %v - want: %v", len(elements), tt.want)
			}
		})
	}
}
