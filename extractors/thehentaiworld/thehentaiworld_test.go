package thehentaiworld

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/config"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		/*{
			Name: "Single video",
			URL:  "https://thehentaiworld.com/videos/shiranui-mai-akiyamaryo-king-of-fighters-5/",
			Want: 1,
		}, {
			Name: "Single 3d-cgi-hentai-images",
			URL:  "https://thehentaiworld.com/3d-cgi-hentai-images/victoria-chase-ubermachine-life-is-strange/",
			Want: 24,
		}, {
			Name: "Single Gallery gif-animated-hentai-images",
			URL:  "https://thehentaiworld.com/gif-animated-hentai-images/hentai-gif-10/",
			Want: 1,
		}, {
			Name: "Single Gallery hentai-cosplay-images",
			URL:  "https://thehentaiworld.com/hentai-cosplay-images/utsushimi-camie-wowmalpal-my-hero-academia/",
			Want: 1,
		}, {
			Name: "Single Gallery hentai-cosplay-images",
			URL:  "https://thehentaiworld.com/hentai-doujinshi/the-start-of-a-harem-juna-juna-juice-my-hero-academia/",
			Want: 1,
		}, {
			Name: "Single Gallery hentai-images",
			URL:  "https://thehentaiworld.com/hentai-images/nico-robin-tit-flash-one-piece-hentai-image/",
			Want: 1,
		}, {
			Name: "Overview",
			URL:  "https://thehentaiworld.com/?new",
			Want: 24,
		}, {
			Name: "Specific Page",
			URL:  "https://thehentaiworld.com/page/4/?s=cyberpunk",
			Want: 24,
		}, */{
			Name: "Tag",
			URL:  "https://thehentaiworld.com/tag/cyberpunk-2077/page/4/",
			Want: 24,
		}, {
			Name: "Mass",
			URL:  "https://thehentaiworld.com/?s=ahri",
			Want: 30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.Name == "Mass" {
				config.Amount = 30
			}

			URLs := parseURL(tt.URL)
			if len(URLs) > tt.Want || len(URLs) == 0 {
				t.Errorf("Got: %v - Want: %v", len(URLs), tt.Want)
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
			Name: "Single video",
			URL:  "hhttps://thehentaiworld.com/videos/ahri-bewyx-league-of-legends-9/",
			Want: 1,
		}, {
			// downloads all images if the post is a image set
			Name: "Single Gallery hentai-images",
			URL:  "https://thehentaiworld.com/hentai-cosplay-images/ahri-helly-von-valentine-league-of-legends-2/",
			Want: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.URL)
			if err != nil {
				t.Error(err)
			}
			if len(data) > tt.Want || len(data) == 0 {
				t.Errorf("Got: %v - Want: %v", len(data), tt.Want)
			}
		})
	}
}
