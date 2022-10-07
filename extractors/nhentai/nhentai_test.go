package nhentai

/*
// make sure that User-Agent makes the browser that did the CF challenge
const userHeader = `cookie: cf_clearance=k2TGEnkzhz_PtHs09vMryROlD4O3UZhrDFrU4svgjdM-1665105987-0-150; csrftoken=bLiwSENr0mqSZZ27wan1xdjLazVFoXnnABJu7DtrhbNRUacpbEZhV0Eggc5lD8m5
user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36`

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
	config.UserHeaders = userHeader
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
		Name string
		Args test.Args
	}{
		{
			Name: "Complete extraction of a doujinshi",
			Args: test.Args{
				URL:     "https://nhentai.net/g/422956/",
				Title:   "Isekai Shoukan IIsan no Tomodachi wa Suki desu ka?",
				Quality: "",
				Size:    0,
			},
		},
		{
			Name: "One page extraction",
			Args: test.Args{
				URL:     "https://nhentai.net/g/297280/14/",
				Title:   "Koe Dashicha Barechau kara!",
				Quality: "",
				Size:    0,
			},
		},
	}
	config.UserHeaders = userHeader
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.Args.URL)
			test.CheckError(t, err)
			test.Check(t, tt.Args, data[0])
		})
	}
}*/
