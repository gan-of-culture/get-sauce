package underhentai

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		URL  string
		want int //number of links
	}{
		{
			name: "Releases 2019",
			URL:  "https://www.underhentai.net/releases-2019/",
			want: 10,
		}, {
			name: "Tag Ahegao",
			URL:  "https://www.underhentai.net/tag/ahegao/",
			want: 6,
		}, {
			name: "Index B",
			URL:  "https://www.underhentai.net/index/b/",
			want: 6,
		}, {
			name: "Normal link",
			URL:  "https://www.underhentai.net/kiss-hug/",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elements, err := ParseURL(tt.URL)
			if err != nil {
				t.Error(err)
			}

			if len(elements) < tt.want {
				t.Errorf("Got: %v - want: %v", len(elements), tt.want)
			}
		})
	}
}

func TestExtractData(t *testing.T) {
	type want struct {
		Title     []string
		StreamLen int
	}
	tests := []struct {
		name string
		URL  string
		want want
	}{
		{
			name: "Single Episode",
			URL:  "https://www.underhentai.net/kiss-hug/",
			want: want{
				Title:     []string{"kiss-hug episode 01"},
				StreamLen: 2,
			},
		}, {
			name: "Multiple Episodes",
			URL:  "https://www.underhentai.net/ochi-mono-rpg-seikishi-luvilias/",
			want: want{
				Title:     []string{"ochi-mono-rpg-seikishi-luvilias episode 01", "ochi-mono-rpg-seikishi-luvilias episode 02", "ochi-mono-rpg-seikishi-luvilias episode 03", "ochi-mono-rpg-seikishi-luvilias episode 04"},
				StreamLen: 8,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := extractData(tt.URL)
			if err != nil {
				t.Error(err)
			}
			for idx, d := range data {
				if d.Title != tt.want.Title[idx] {
					t.Errorf("Got: %v - want: %v", d.Title, tt.want.Title[idx])
				}
				if len(d.Streams) != tt.want.StreamLen {
					t.Errorf("Got: %v - want: %v", len(d.Streams), tt.want.StreamLen)
				}
			}
		})
	}
}
