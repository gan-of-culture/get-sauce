// for image viewing limits please checkout this article: https://github.com/8qwe24657913/E-Hentai-Downloader-NW.js/wiki/E%E2%88%92Hentai-Image-Viewing-Limits
package ehentai

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/test"
)

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
		Args test.Args
	}{
		{
			Name: "Default extraction",
			Args: test.Args{
				URL:     "https://e-hentai.org/g/1559777/dc952bd4c1/",
				Title:   "[Shimashima-PNT (Punita)] Diablo no Shoyuubutsu dakara Suki ni Shite mo Ii yo... | Диабло, мы твои рабы, так что ты можешь делать всё, что пожелаешь… [Russian] [R9N9Ga7] [Digital] - 1",
				Quality: "1280 x 1807",
				Size:    497000,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.Args.URL)
			test.CheckError(t, err)
			test.Check(t, tt.Args, data[0])
		})
	}
}
