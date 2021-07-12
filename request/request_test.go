package request

import (
	"testing"

	"github.com/gan-of-culture/go-hentai-scraper/config"
)

func TestSize(t *testing.T) {
	config.ShowInfo = true
	t.Run("Default test", func(t *testing.T) {
		size, err := Size("https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png", "")
		if err != nil {
			t.Error(err)
		}

		if size == 0 {
			t.Errorf("Got: %v - want: %v", size, "more than 0 Bytes")
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
			t.Errorf("Got: %v - want: %v", htmlString, "a string")
		}
	})
}
