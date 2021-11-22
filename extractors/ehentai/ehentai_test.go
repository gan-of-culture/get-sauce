// for image viewing limits please checkout this article: https://github.com/8qwe24657913/E-Hentai-Downloader-NW.js/wiki/E%E2%88%92Hentai-Image-Viewing-Limits
package ehentai

import "testing"

func TestPaseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		//expect atlest ... galleries
		numberOfGalleries int
	}{
		{
			Name:              "Parse page of galleries",
			URL:               "https://e-hentai.org/?page=1&f_cats=1021",
			numberOfGalleries: 5,
		}, {
			Name:              "Single gallery",
			URL:               "https://e-hentai.org/g/1559777/dc952bd4c1/",
			numberOfGalleries: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			URLs := parseURL(tt.URL)
			if len(URLs) < tt.numberOfGalleries {
				t.Errorf("Got: %v - Want: %v", len(URLs), tt.numberOfGalleries)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		// expect atleast ... data structs
		numberOfData int
	}{
		{
			Name:         "Single gallery",
			URL:          "https://e-hentai.org/g/1559777/dc952bd4c1/",
			numberOfData: 26,
		},
		//commented out because of performance - look at the parseURl test instead
		/*{
			Name:         "Parse page of galleries",
			URL:          "https://e-hentai.org/?page=1&f_cats=1021",
			numberOfData: 150,
		},*/
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.URL)
			if err != nil {
				t.Error(err)
			}
			if len(data) < tt.numberOfData {
				t.Errorf("Got: %v - Want: %v", len(data), tt.numberOfData)
			}
		})
	}
}
