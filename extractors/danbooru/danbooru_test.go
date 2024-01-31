package danbooru

import (
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
			Name: "Overview page",
			URL:  "https://danbooru.donmai.us/posts?page=3&tags=fire_emblem",
			Want: 2,
		}, {
			Name: "Example Post",
			URL:  "https://danbooru.donmai.us/posts/3749687",
			Want: 1,
		},
	}
	config.FakeHeaders["User-Agent"] = "Mozilla/4.0 (compatible; MSIE 9.0; Windows NT 6.1)"
	defer func() {
		config.FakeHeaders["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.81 Safari/537.36"
	}()
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			URLs, err := parseURL(tt.URL)
			test.CheckError(t, err)
			if len(URLs) < tt.Want {
				t.Errorf("Got: %v - Want: %v", len(URLs), tt.Want)
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
				Size:    0,
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
