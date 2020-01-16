package nhentai

import "testing"

func TestParseURL(t *testing.T) {
	type want struct {
		magicNumer string
		page       string
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Only magic number supplied",
			url:  "https://nhentai.net/g/297495/",
			want: want{
				magicNumer: "297495",
				page:       "",
			},
		}, {
			name: "magic number and page number supplied",
			url:  "https://nhentai.net/g/297485/9/",
			want: want{
				magicNumer: "297485",
				page:       "9",
			},
		}, {
			name: "Incorrect url",
			url:  "https://nhentai.net/g/",
			want: want{
				magicNumer: "",
				page:       "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			magicNumer, page := ParseURL(tt.url)

			if magicNumer != tt.want.magicNumer {
				t.Errorf("Got: %v - want: %v", magicNumer, tt.want)
			}

			if page != tt.want.page {
				t.Errorf("Got: %v - want: %v", magicNumer, tt.want)
			}
		})
	}
}

func TestExtract(t *testing.T) {

	tests := []struct {
		name  string
		url   string
		title string
		want  int
	}{
		{
			name:  "Complete extraction of a doujinshi",
			url:   "https://nhentai.net/g/297485/",
			title: "Isekai Shoukan II - Elf na Onee-san no Tomodachi wa Suki desu ka?",
			want:  43,
		}, {
			name:  "One page extraction",
			url:   "https://nhentai.net/g/297280/14/",
			title: "(C97) [H@BREAK (Itose Ikuto)] Koe Dashicha Barechau kara! [English]",
			want:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Extract(tt.url)
			if err != nil {
				t.Error(err)
			}
			dataLen := len(data)
			if dataLen != tt.want {
				t.Errorf("Got: %v - want: %v", dataLen, tt.want)
			}
		})
	}
}
