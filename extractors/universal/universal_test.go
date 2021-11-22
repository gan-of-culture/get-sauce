package universal

import "testing"

func TestExtract(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		site string
		Want int
	}{
		{
			Name: "Imigur",
			URL:  "http://i.imgur.com/06YTAjg.jpg",
			site: "imigur",
			Want: 1,
		}, {
			Name: "awwni",
			URL:  "http://cdn.awwni.me/16c8v.jpg",
			site: "awwni",
			Want: 1,
		}, {
			Name: "with bloat after ext",
			URL:  "https://img.rule34.xxx//images/1979/b84be533024a3d1dcc6b01c0cb7358c9.jpeg?2686173",
			site: "rule34.xxx",
			Want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.URL)
			if err != nil {
				t.Error(err)
			}
			if len(data) != tt.Want {
				t.Errorf("Got: %v - Want: %v", len(data), tt.Want)
			}
		})
	}
}
