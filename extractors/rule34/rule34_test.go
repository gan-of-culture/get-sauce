package rule34

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/static"
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
	type Want struct {
		Title   string
		Type    static.DataType
		DataLen int
	}
	tests := []struct {
		Name string
		URL  string
		Want Want
	}{
		{
			Name: "Test image",
			URL:  "https://rule34.paheal.net/post/view/3427635",
			Want: Want{
				Title:   "Magical_Sempai_(Series) Magician_Sempai skyfreedom 3427635",
				Type:    static.DataTypeImage,
				DataLen: 1,
			},
		}, {
			Name: "Test video",
			URL:  "https://rule34.paheal.net/post/view/3464181",
			Want: Want{
				Title:   "Hv54rDSL Nier Nier_Automata YoRHa_No.2_Type_B animated audiodude blender sound webm 3464181",
				Type:    static.DataTypeVideo,
				DataLen: 1,
			},
		}, {
			Name: "Test gif",
			URL:  "https://rule34.paheal.net/post/view/3461411",
			Want: Want{
				Title:   "World_of_Warcraft animated blood_elf 3461411",
				Type:    static.DataTypeImage,
				DataLen: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			config.Amount = 26
			elements, err := New().Extract(tt.URL)
			if err != nil {
				t.Error("elements has error or is too big for single tests")
			}
			act := Want{
				Title:   elements[0].Title,
				Type:    elements[0].Type,
				DataLen: len(elements),
			}
			if act != tt.Want {
				t.Errorf("Got: %v - Want: %v", act, tt.Want)
			}
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
				t.Error("elements has error or is too big for single tests")
			}
			if len(elements) != tt.Want {
				t.Errorf("Got: %v - Want: %v", len(elements), tt.Want)
			}
		})
	}
}
