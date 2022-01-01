package imgboard

import (
	"strings"
	"testing"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/test"
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
	tests := []struct {
		Name string
		Args test.Args
	}{
		{
			Name: "Test image",
			Args: test.Args{
				URL:     "https://gelbooru.com/index.php?page=post&s=view&id=5888807",
				Title:   "gelbooru_5888807",
				Quality: "800x1280",
				Size:    679127,
			},
		},
		{
			Name: "Test video",
			Args: test.Args{
				URL:     "https://rule34.xxx/index.php?page=post&s=view&id=4470579",
				Title:   "rule34_4470579",
				Quality: "1280x720",
				Size:    4134392,
			},
		},
		{
			Name: "Test image konachan",
			Args: test.Args{
				URL:     "https://konachan.com/post/show/323560/anthropomorphism-azur_lane-black_hair-blush-breast",
				Title:   "konachan_323560",
				Quality: "1371x1029",
				Size:    866039,
			},
		},
		{
			Name: "Test image yande.re",
			Args: test.Args{
				URL:     "https://yande.re/post/show/745150",
				Title:   "yande_745150",
				Quality: "2275x3660",
				Size:    4682030,
			},
		},
		{
			Name: "Test image tbib",
			Args: test.Args{
				URL:     "https://tbib.org/index.php?page=post&s=view&id=9022091",
				Title:   "tbib_9022091",
				Quality: "2869x5100",
				Size:    2179461,
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
			test.CheckError(t, err)
			if len(elements) < tt.Want {
				t.Errorf("Got: %v - Want: %v", len(elements), tt.Want)
			}
		})
	}
}
