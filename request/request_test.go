package request

import "testing"

func TestSize(t *testing.T) {
	t.Run("Default test", func(t *testing.T) {
		size, err := Size("https://wikipedia.de/img/Wikipedia-logo-v2-de.svg", "")
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
