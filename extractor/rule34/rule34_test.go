package rule34

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elements := ParseURL(tt.url)

			if len(elements) < tt.want {
				t.Errorf("Got: %v - want: %v", len(elements), tt.want)
			}
		})
	}

}

func TestExtract(t *testing.T) {
	type want struct {
		Title     string
		Type      string
		StreamLen int
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
				Title:     "Magical_Sempai_(Series) Magician_Sempai skyfreedom",
				Type:      "image",
				StreamLen: 1,
			},
		}, {
			name: "Test video",
			url:  "https://rule34.paheal.net/post/view/3464181",
			want: want{
				Title:     "Hv54rDSL Nier Nier_Automata YoRHa_No.2_Type_B animated audiodude blender sound webm",
				Type:      "video",
				StreamLen: 1,
			},
		}, {
			name: "Test gif",
			url:  "https://rule34.paheal.net/post/view/3461411",
			want: want{
				Title:     "World_of_Warcraft animated blood_elf",
				Type:      "gif",
				StreamLen: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elements, err := Extractor(tt.url)
			if err != nil {
				t.Error("elements has error or is too big for single tests")
			}
			act := want{
				Title:     elements[0].Title,
				Type:      elements[0].Type,
				StreamLen: len(elements[0].Streams),
			}
			if act != tt.want {
				t.Errorf("Got: %v - want: %v", act, tt.want)
			}
		})
	}
}
