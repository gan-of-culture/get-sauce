package request

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/static"
)

func TestSize(t *testing.T) {
	config.ShowInfo = true
	t.Run("Default test", func(t *testing.T) {
		size, err := Size("https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png", "")
		if err != nil {
			t.Error(err)
		}

		if size == 0 {
			t.Errorf("Got: %v - Want: %v", size, "more than 0 Bytes")
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("Default test", func(t *testing.T) {
		htmlString, err := Get("https://github.com/")
		if err != nil {
			t.Error(err)
		}

		if htmlString == "" {
			t.Errorf("Got: %v - Want: %v", htmlString, "a string")
		}
	})
}

func TestPost(t *testing.T) {
	t.Run("Default test", func(t *testing.T) {
		data, err := PostAsBytesWithHeaders("https://www.google.com/", map[string]string{"Referer": "https://google.com"})
		if err != nil {
			t.Error(err)
		}

		if len(data) < 1 {
			t.Errorf("Got: %v - Want: %v", data, "some bytes")
		}
	})
}

func TestGetWReferer(t *testing.T) {
	t.Run("Default test", func(t *testing.T) {
		htmlString, err := GetWithHeaders("https://github.com/", map[string]string{
			"referer": "https://github.com/",
		})
		if err != nil {
			t.Error(err)
		}

		if htmlString == "" {
			t.Errorf("Got: %v - Want: %v", htmlString, "a string")
		}
	})
}

func TestExtractHLS(t *testing.T) {
	tests := []struct {
		Name    string
		URL     string
		Headers map[string]string
		Want    map[string]*static.Stream
	}{
		{
			Name: "HLS where stream order is from small to high",
			URL:  "https://hentaidoge.org/hls/K/kuroinu-ii-animation/1/playlist.m3u8",
			Headers: map[string]string{
				"Referer": "https://hentaidoge.org/hls/K/kuroinu-ii-animation/1/playlist.m3u8",
			},
			Want: map[string]*static.Stream{
				"0": {
					Type:    static.DataTypeVideo,
					Quality: "1280x720",
				},
				"1": {
					Type:    static.DataTypeVideo,
					Quality: "842x480",
				},
				"2": {
					Type:    static.DataTypeVideo,
					Quality: "640x360",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			streams, err := ExtractHLS(tt.URL, tt.Headers)
			if err != nil {
				t.Error(err)
			}
			for k, v := range streams {
				if v.Quality != tt.Want[k].Quality {
					t.Errorf("Got: %v - Want: %v", v.Quality, tt.Want[k].Quality)
				}
			}
		})
	}
}
