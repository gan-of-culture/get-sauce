package imgboard

import (
	"strings"
	"testing"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/static"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single Page",
			URL:  "https://rule34.xxx/index.php?page=post&s=list&tags=world_of_warcraft+video+draenei&pid=378",
			// atleast more than 2
			Want: 2,
		}, {
			Name: "Mass extract page booru project",
			URL:  "https://tbib.org/index.php?page=post&s=list&tags=1girl+solo+uncensored+full_body+&pid=0",
			Want: 100,
		},
		{
			Name: "Single image booru project",
			URL:  "https://rule34.xxx/index.php?page=post&s=view&id=4470590",
			Want: 1,
		}, {
			Name: "Overview page booru project",
			URL:  "https://rule34.xxx/index.php?page=post&s=list&tags=all",
			// atleast more than 2
			Want: 2,
		}, {
			Name: "Mass extract page booru project",
			URL:  "https://rule34.xxx/index.php?page=post&s=list&tags=all",
			// atleast more than 2
			Want: 100,
		}, {
			Name: "Mass extract page booru project2",
			URL:  "https://gelbooru.com/index.php?page=post&s=list&tags=all",
			// atleast more than 2
			Want: 100,
		}, {
			Name: "Single image",
			URL:  "https://yande.re/post/show/745150",
			Want: 1,
		}, {
			Name: "Overview page booru project",
			URL:  "https://konachan.com/post?tags=uncensored",
			// atleast more than 2
			Want: 2,
		}, {
			Name: "Mass extract page",
			URL:  "https://konachan.com/post?tags=uncensored",
			// atleast more than 2
			Want: 100,
		}, {
			Name: "Mass extract page2",
			URL:  "https://yande.re/post?tags=tateha&page=2",
			// atleast more than 2
			Want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if strings.HasPrefix(tt.Name, "Mass extract") {
				config.Amount = 101
			}
			elements := parseURL(tt.URL)

			if len(elements) < tt.Want {
				t.Errorf("Got: %v - Want: %v", len(elements), tt.Want)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	type Want struct {
		Title   string
		Type    static.DataType
		DataLen int
	}
	tests := []struct {
		Name string
		URL  string
		Want Want
	}{
		{
			Name: "Test image",
			URL:  "https://gelbooru.com/index.php?page=post&s=view&id=5888807",
			Want: Want{
				Title:   "gelbooru_5888807",
				Type:    static.DataTypeImage,
				DataLen: 1,
			},
		}, {
			Name: "Test video",
			URL:  "https://rule34.xxx/index.php?page=post&s=view&id=4470579",
			Want: Want{
				Title:   "rule34_4470579",
				Type:    static.DataTypeVideo,
				DataLen: 1,
			},
		}, {
			Name: "Test image konachan",
			URL:  "https://konachan.com/post/show/323560/anthropomorphism-azur_lane-black_hair-blush-breast",
			Want: Want{
				Title:   "konachan_323560",
				Type:    static.DataTypeImage,
				DataLen: 1,
			},
		}, {
			Name: "Test image yande.re",
			URL:  "https://yande.re/post/show/745150",
			Want: Want{
				Title:   "yande_745150",
				Type:    static.DataTypeImage,
				DataLen: 1,
			},
		}, {
			Name: "Test image tbib",
			URL:  "https://tbib.org/index.php?page=post&s=view&id=9022091",
			Want: Want{
				Title:   "tbib_9022091",
				Type:    static.DataTypeImage,
				DataLen: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			elements, err := New().Extract(tt.URL)
			if err != nil {
				t.Error(err)
			}
			act := Want{
				Title:   elements[0].Title,
				Type:    elements[0].Type,
				DataLen: len(elements),
			}
			if act != tt.Want {
				t.Errorf("Got: %v - Want: %v", act, tt.Want)
			}
		})
	}
}

func TestExtractDataFromDirectLink(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "DL extract yande.re",
			URL:  "https://yande.re/post?",
			Want: 10,
		}, {
			Name: "DL extract konachan",
			URL:  "https://konachan.com/post?tags=",
			Want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			config.Amount = 10
			elements, err := New().Extract(tt.URL)
			if err != nil {
				t.Error(err)
			}
			if len(elements) < tt.Want {
				t.Errorf("Got: %v - Want: %v", len(elements), tt.Want)
			}
		})
	}
}
