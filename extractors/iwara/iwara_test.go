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
			URL:  "https://iwara.tv/video/1ye3vfv2bpfmpe0k2",
			Want: 1,
		}, {
			Name: "Single images",
			URL:  "https://iwara.tv/image/x6hVrNaf0WVdLE/nico-tomoare-provocation-dance-preview-mmdd",
			Want: 1,
		}, {
			Name: "Mass",
			URL:  "https://iwara.tv/images?sort=date&rating=ecchi&page=1",
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
				URL:     "https://www.iwara.tv/video/zwkobu5eevi17m309/14",
				Title:   "14.舞夜の侍奉Ⅱ💕",
				Quality: "Source",
				Size:    90608128,
			},
		},
		{
			Name: "Single images post",
			Args: test.Args{
				URL:     "https://iwara.tv/image/x6hVrNaf0WVdLE/nico-tomoare-provocation-dance-preview-mmdd",
				Title:   "【コイカツ】 Nico Thick Tomoare Provocation Dance Preview 【MMDD】",
				Quality: "1920x1080",
				Size:    454545,
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
