package rule34

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/static"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		amount int
		want   int
	}{
		{
			name: "Single image",
			url:  "https://rule34.paheal.net/post/view/3464197",
			want: 1,
		}, {
			name: "Single video",
			url:  "https://rule34.paheal.net/post/view/3464181",
			want: 1,
		}, {
			name: "Overview page",
			url:  "https://rule34.paheal.net/post/list/2",
			// atleast more than 2
			want: 2,
		}, {
			name: "Mass extract page",
			url:  "https://rule34.paheal.net/post/list/1",
			// atleast more than 2
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Mass extract page" {
				config.Amount = 101
			}
			elements := parseURL(tt.url)

			if len(elements) < tt.want {
				t.Errorf("Got: %v - want: %v", len(elements), tt.want)
			}
		})
	}

}

func TestExtract(t *testing.T) {
	type want struct {
		Title   string
		Type    static.DataType
		DataLen int
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Test image",
			url:  "https://rule34.paheal.net/post/view/3427635",
			want: want{
				Title:   "Magical_Sempai_(Series) Magician_Sempai skyfreedom 3427635",
				Type:    static.DataTypeImage,
				DataLen: 1,
			},
		}, {
			name: "Test video",
			url:  "https://rule34.paheal.net/post/view/3464181",
			want: want{
				Title:   "Hv54rDSL Nier Nier_Automata YoRHa_No.2_Type_B animated audiodude blender sound webm 3464181",
				Type:    static.DataTypeVideo,
				DataLen: 1,
			},
		}, {
			name: "Test gif",
			url:  "https://rule34.paheal.net/post/view/3461411",
			want: want{
				Title:   "World_of_Warcraft animated blood_elf 3461411",
				Type:    static.DataTypeImage,
				DataLen: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Amount = 26
			elements, err := New().Extract(tt.url)
			if err != nil {
				t.Error("elements has error or is too big for single tests")
			}
			act := want{
				Title:   elements[0].Title,
				Type:    elements[0].Type,
				DataLen: len(elements),
			}
			if act != tt.want {
				t.Errorf("Got: %v - want: %v", act, tt.want)
			}
		})
	}
}

func TestMassExtract(t *testing.T) {

	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Test mass",
			url:  "https://rule34.paheal.net/post/list/1",
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Amount = 10
			elements, err := New().Extract(tt.url)
			if err != nil {
				t.Error("elements has error or is too big for single tests")
			}
			if len(elements) != tt.want {
				t.Errorf("Got: %v - want: %v", len(elements), tt.want)
			}
		})
	}
}
