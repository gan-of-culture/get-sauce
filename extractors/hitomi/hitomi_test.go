package hitomi

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/test"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single Gallery",
			URL:  "https://hitomi.la/doujinshi/%E3%82%B8%E3%82%A7%E3%83%B3%E3%83%88%E3%83%AB%E3%82%B3%E3%83%8D%E3%82%AF%E3%83%88!re:dive-%E6%97%A5%E6%9C%AC%E8%AA%9E-1905632.html",
			Want: 1,
		}, {
			Name: "Tag",
			URL:  "https://hitomi.la/tag/artbook-all.html",
			Want: 25,
		}, {
			Name: "Tag different page",
			URL:  "https://hitomi.la/tag/artbook-all.html?page=2",
			Want: 25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			URLs := parseURL(tt.URL)
			if len(URLs) < tt.Want || len(URLs) == 0 {
				t.Errorf("Got: %v - Want: %v", len(URLs), tt.Want)
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
			Name: "Single Doujin",
			Args: test.Args{
				URL:     "https://hitomi.la/doujinshi/m-o-muke-onaclu-_shinjin-kenshuu-hen-english-2102221.html",
				Title:   "M-o Muke OnaClu _Shinjin Kenshuu Hen",
				Quality: "",
				Size:    0,
			},
		},
		{
			Name: "Single Manga",
			Args: test.Args{
				URL:     "https://hitomi.la/manga/%E7%8C%A5%E8%A4%BB%E3%83%9F%E3%82%B5%E3%82%A4%E3%83%AB-%E6%97%A5%E6%9C%AC%E8%AA%9E-440479.html",
				Title:   "Waisetsu Missile",
				Quality: "",
				Size:    0,
			},
		},
		{
			Name: "Single Game CG",
			Args: test.Args{
				URL:     "https://hitomi.la/gamecg/taimanin-rpgx-2194813.html#1",
				Title:   "Taimanin RPGX",
				Quality: "",
				Size:    0,
			},
		},
		{
			Name: "Single Artist CG",
			Args: test.Args{
				URL:     "https://hitomi.la/cg/%E3%82%A6%E3%83%81%E3%81%AE%E5%BC%9F%E3%83%9E%E3%82%B8%E3%81%A7%E3%83%87%E3%82%AB%E3%82%A4%E3%82%93%E3%81%A0%E3%81%91%E3%81%A9%E8%A6%8B%E3%81%AB%E3%81%93%E3%81%AA%E3%81%84%EF%BC%9F-%E6%97%A5%E6%9C%AC%E8%AA%9E-1258357.html#1",
				Title:   "Uchi no Otouto Maji de Dekain dakedo Mi ni Konai?",
				Quality: "",
				Size:    0,
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
