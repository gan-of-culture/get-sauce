package nhentai

import "testing"

func TestParseURL(t *testing.T) {
	type Want struct {
		magicNumbers int
		page         string
	}
	tests := []struct {
		Name string
		URL  string
		Want Want
	}{
		{
			Name: "Only magic number supplied",
			URL:  "https://nhentai.net/g/297495/",
			Want: Want{
				magicNumbers: 1,
				page:         "",
			},
		}, {
			Name: "magic number and page number supplied",
			URL:  "https://nhentai.net/g/297485/9/",
			Want: Want{
				magicNumbers: 1,
				page:         "9",
			},
		}, {
			Name: "Incorrect url",
			URL:  "https://nhentai.net/g/",
			Want: Want{
				magicNumbers: 0,
				page:         "",
			},
		}, {
			Name: "Doujin collection",
			URL:  "https://nhentai.net/search/?q=dragon",
			Want: Want{
				magicNumbers: 22,
				page:         "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			magicNumber, page := parseURL(tt.URL)
			if len(magicNumber) < tt.Want.magicNumbers {
				t.Errorf("Got: %v - Want: %v", len(magicNumber), tt.Want)
			}

			if page != tt.Want.page {
				t.Errorf("Got: %v - Want: %v", len(magicNumber), tt.Want)
			}
		})
	}
}

func TestExtract(t *testing.T) {

	tests := []struct {
		Name  string
		URL   string
		title string
		Want  int
	}{
		{
			Name:  "Complete extraction of a doujinshi",
			URL:   "https://nhentai.net/g/297485/",
			title: "Isekai Shoukan II - Elf na Onee-san no Tomodachi wa Suki desu ka?",
			Want:  43,
		}, {
			Name:  "One page extraction",
			URL:   "https://nhentai.net/g/297280/14/",
			title: "(C97) [H@BREAK (Itose Ikuto)] Koe Dashicha Barechau kara! [English]",
			Want:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.URL)
			if err != nil {
				t.Error(err)
			}
			URLlen := len(data[0].Streams["0"].URLs)
			if URLlen != tt.Want {
				t.Errorf("Got: %v - Want: %v", URLlen, tt.Want)
			}
		})
	}
}
