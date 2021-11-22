package hentaihaven

/*
import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single Episode",
			URL:  "https://hentaihaven.xxx/watch/ero-konbini-tenchou/episode-4/",
			Want: 1,
		}, {
			Name: "Series",
			URL:  "https://hentaihaven.xxx/watch/ero-konbini-tenchou/",
			Want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
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
			Name: "Single Episode",
			URL:  "https://hentaihaven.xxx/watch/showtime-uta-no-onee-san-datte-shitai/episode-3/",
			Want: 1,
		}, {
			Name: "Series",
			URL:  "https://hentaihaven.xxx/watch/ero-konbini-tenchou/",
			Want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.URL)
			if err != nil {
				t.Error(err)
			}
			if len(data) > tt.Want {
				t.Errorf("Got: %v - Want: %v", len(data), tt.Want)
			}
		})
	}
}*/
