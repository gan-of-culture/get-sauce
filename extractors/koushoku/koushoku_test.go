package koushoku

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
			URL:  "https://koushoku.org/archive/8411/intercourse-inn",
			Want: 1,
		}, {
			Name: "Single Page",
			URL:  "https://koushoku.org/archive/8411/intercourse-inn/10",
			Want: 1,
		}, {
			Name: "Tag",
			URL:  "https://koushoku.org/tags/box-set",
			Want: 25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
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
			Name: "Single Gallery",
			Args: test.Args{
				URL:     "https://koushoku.org/archive/7915/no-virgins-allowed-the-time-a-creepy-otaku-like-me-helped-the-class-gyarus-lose-their-virginity",
				Title:   "No Virgins Allowed - The Time a Creepy Otaku Like Me Helped the Class Gyarus Lose Their Virginity",
				Quality: "",
				Size:    99300000,
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
