package ehentai

import "testing"

func TestPaseURL(t *testing.T) {
	tests := []struct {
		name string
		URL  string
		//expect atlest ... galleries
		numberOfGalleries int
	}{
		{
			name:              "Parse page of galleries",
			URL:               "https://e-hentai.org/?page=1&f_cats=1021",
			numberOfGalleries: 5,
		}, {
			name:              "Single gallery",
			URL:               "https://e-hentai.org/g/1559777/dc952bd4c1/",
			numberOfGalleries: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			URLs := parseURL(tt.URL)
			if len(URLs) < tt.numberOfGalleries {
				t.Errorf("Got: %v - want: %v", len(URLs), tt.numberOfGalleries)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		name string
		URL  string
		// expect atleast ... data structs
		numberOfData int
	}{
		{
			name:         "Single gallery",
			URL:          "https://e-hentai.org/g/1559777/dc952bd4c1/",
			numberOfData: 26,
		},
		//commented out because of performance - look at the parseURl test instead
		/*{
			name:         "Parse page of galleries",
			URL:          "https://e-hentai.org/?page=1&f_cats=1021",
			numberOfData: 150,
		},*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := New().Extract(tt.URL)
			if err != nil {
				t.Error(err)
			}
			if len(data) < tt.numberOfData {
				t.Errorf("Got: %v - want: %v", len(data), tt.numberOfData)
			}
		})
	}
}
