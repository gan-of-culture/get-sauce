package manhwa18

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single Episode",
			URL:  "https://manhwa18.tv/taimanin-shiranui-inyoku-no-dorei-shoufu-episode-1",
			Want: 1,
		}, {
			Name: "Mass",
			URL:  "https://manhwa18.tv/uncensored",
			Want: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			URLs := parseURL(tt.URL)
			if len(URLs) < tt.Want {
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
			Name: "Single Episode",
			URL:  "https://manhwa18.tv/kazoku-haha-to-shimai-no-kyousei-episode-3",
			Want: 1,
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
