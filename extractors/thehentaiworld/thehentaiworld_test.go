package thehentaiworld

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/v2/config"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		/*{
			name: "Single video",
			url:  "https://thehentaiworld.com/videos/shiranui-mai-akiyamaryo-king-of-fighters-5/",
			want: 1,
		}, {
			name: "Single 3d-cgi-hentai-images",
			url:  "https://thehentaiworld.com/3d-cgi-hentai-images/victoria-chase-ubermachine-life-is-strange/",
			want: 24,
		}, {
			name: "Single Gallery gif-animated-hentai-images",
			url:  "https://thehentaiworld.com/gif-animated-hentai-images/hentai-gif-10/",
			want: 1,
		}, {
			name: "Single Gallery hentai-cosplay-images",
			url:  "https://thehentaiworld.com/hentai-cosplay-images/utsushimi-camie-wowmalpal-my-hero-academia/",
			want: 1,
		}, {
			name: "Single Gallery hentai-cosplay-images",
			url:  "https://thehentaiworld.com/hentai-doujinshi/the-start-of-a-harem-juna-juna-juice-my-hero-academia/",
			want: 1,
		}, {
			name: "Single Gallery hentai-images",
			url:  "https://thehentaiworld.com/hentai-images/nico-robin-tit-flash-one-piece-hentai-image/",
			want: 1,
		}, {
			name: "Overview",
			url:  "https://thehentaiworld.com/?new",
			want: 24,
		}, {
			name: "Specific Page",
			url:  "https://thehentaiworld.com/page/4/?s=cyberpunk",
			want: 24,
		}, */{
			name: "Tag",
			url:  "https://thehentaiworld.com/tag/cyberpunk-2077/page/4/",
			want: 24,
		}, {
			name: "Mass",
			url:  "https://thehentaiworld.com/?s=ahri",
			want: 30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Mass" {
				config.Amount = 30
			}

			urls := parseURL(tt.url)
			if len(urls) > tt.want || len(urls) == 0 {
				t.Errorf("Got: %v - want: %v", len(urls), tt.want)
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
			name: "Single video",
			url:  "hhttps://thehentaiworld.com/videos/ahri-bewyx-league-of-legends-9/",
			want: 1,
		}, {
			// downloads all images if the post is a image set
			name: "Single Gallery hentai-images",
			url:  "https://thehentaiworld.com/hentai-cosplay-images/ahri-helly-von-valentine-league-of-legends-2/",
			want: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := New().Extract(tt.url)
			if err != nil {
				t.Error(err)
			}
			if len(data) > tt.want || len(data) == 0 {
				t.Errorf("Got: %v - want: %v", len(data), tt.want)
			}
		})
	}
}
