package rule34

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/test"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name   string
		URL    string
		amount int
		Want   int
	}{
		{
			Name: "Single image",
			URL:  "https://rule34.paheal.net/post/view/3464197",
			Want: 1,
		}, {
			Name: "Single video",
			URL:  "https://rule34.paheal.net/post/view/3464181",
			Want: 1,
		}, {
			Name: "Overview page",
			URL:  "https://rule34.paheal.net/post/list/2",
			// atleast more than 2
			Want: 2,
		}, {
			Name: "Mass extract page",
			URL:  "https://rule34.paheal.net/post/list/1",
			// atleast more than 2
			Want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.Name == "Mass extract page" {
				config.Amount = 101
			}
			elements := parseURL(tt.URL)

			if len(elements) < tt.Want {
				t.Errorf("Got: %v - Want: %v", len(elements), tt.Want)
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
				URL:     "https://rule34.paheal.net/post/view/3464181",
				Title:   "Hv54rDSL Nier_(series) Nier_Automata YoRHa_No.2_Type_B animated audiodude blender sound webm 3464181",
				Quality: "540 x 1280",
				Size:    7503936,
			},
		},
		{
			Name: "Single image",
			Args: test.Args{
				URL:     "https://rule34.paheal.net/post/view/3427635",
				Title:   "Magical_Sempai_(series) Magician_Sempai skyfreedom 3427635",
				Quality: "1800 x 1269",
				Size:    590974,
			},
		},
		{
			Name: "Single GIF",
			Args: test.Args{
				URL:     "https://rule34.paheal.net/post/view/3461411",
				Title:   "World_of_Warcraft animated blood_elf 3461411",
				Quality: "480 x 854",
				Size:    7811055,
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

func TestMassExtract(t *testing.T) {

	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Test mass",
			URL:  "https://rule34.paheal.net/post/list/1",
			Want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			config.Amount = 10
			elements, err := New().Extract(tt.URL)
			if err != nil {
				test.CheckError(t, err)
			}
			if len(elements) != tt.Want {
				t.Errorf("Got: %v - Want: %v", len(elements), tt.Want)
			}
		})
	}
}
