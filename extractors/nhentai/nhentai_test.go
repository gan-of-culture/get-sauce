package nhentai

/*
// make sure that User-Agent makes the browser that did the CF challenge
const userHeader = `Cookie: cf_clearance=pScTgd92T67tqxJkpQNkbuQ6A4cA.auLs9KaKtECOBE-1681481968-0-160; csrftoken=Exvbm0X8mhEXx02g6hzDkEOfFRb841cgmldz9UewkFWgEXX7fDsalXalHglyrCh0
User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/112.0`

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
				Title:   "Onaho Hanbaiin no Oshigoto wa Taihen desu!",
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
