package iwara

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/test"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single video",
			URL:  "https://www.iwara.tv/video/gz8d8ik2zlhoqjk1w/rbq",
			Want: 1,
		}, {
			Name: "Single images",
			URL:  "https://www.iwara.tv/image/WThZZG81z25j4I/hell-apocalypsemash-kyrielight-edition004",
			Want: 1,
		}, {
			Name: "Mass",
			URL:  "https://www.iwara.tv/images?sort=date&rating=ecchi&page=1",
			Want: 40,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.Name == "Mass" {
				config.Amount = 40
			}

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
		Args test.Args
	}{
		{
			Name: "Single video",
			Args: test.Args{
				URL:   "https://www.iwara.tv/video/gz8d8ik2zlhoqjk1w/rbq",
				Title: "抖音风-十位RBQ的联合参演！",
				Size:  185614559,
			},
		},
		{
			Name: "Single images post",
			Args: test.Args{
				URL:     "https://www.iwara.tv/image/x6hVrNaf0WVdLE/nico-tomoare-provocation-dance-preview-mmdd",
				Title:   "【コイカツ】 Nico Tomoare Provocation Dance Preview 【MMDD】",
				Quality: "1280x720",
				Size:    294562,
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
