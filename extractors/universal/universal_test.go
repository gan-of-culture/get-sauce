package universal

import "testing"

func TestExtract(t *testing.T) {
	tests := []struct {
		name string
		url  string
		site string
		want int
	}{
		{
			name: "Imigur",
			url:  "http://i.imgur.com/06YTAjg.jpg",
			site: "imigur",
			want: 1,
		}, {
			name: "awwni",
			url:  "http://cdn.awwni.me/16c8v.jpg",
			site: "awwni",
			want: 1,
		}, {
			name: "with bloat after ext",
			url:  "https://img.rule34.xxx//images/1979/b84be533024a3d1dcc6b01c0cb7358c9.jpeg?2686173",
			site: "rule34.xxx",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := New().Extract(tt.url)
			if err != nil {
				t.Error(err)
			}
			if len(data) != tt.want {
				t.Errorf("Got: %v - want: %v", len(data), tt.want)
			}
		})
	}
}
