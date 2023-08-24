package request

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/test"
)

func TestSize(t *testing.T) {
	config.ShowInfo = true
	t.Run("Default test", func(t *testing.T) {
		size, err := Size("https://www.google.com/images/branding/googlelogo/2x/googlelogo_color_272x92dp.png", "")
		test.CheckError(t, err)

		if size == 0 {
			t.Errorf("Got: %v - Want: %v", size, "more than 0 Bytes")
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("Default test", func(t *testing.T) {
		htmlString, err := Get("https://github.com/")
		test.CheckError(t, err)

		if htmlString == "" {
			t.Errorf("Got: %v - Want: %v", htmlString, "a string")
		}
	})
}

func TestPost(t *testing.T) {
	t.Run("Default test", func(t *testing.T) {
		data, err := PostAsBytesWithHeaders("https://www.google.com/", map[string]string{"Referer": "https://google.com"}, nil)
		test.CheckError(t, err)

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
		test.CheckError(t, err)

		if htmlString == "" {
			t.Errorf("Got: %v - Want: %v", htmlString, "a string")
		}
	})
}
