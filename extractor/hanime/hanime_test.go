package hanime

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Page 1",
			url:  "https://hanime.tv/browse/images?page=1",
			want: 2,
		}, {
			name: "Image",
			url:  "https://htvassets.club/uploads/776000/776272.png",
			want: 1,
		}, {
			name: "Gif",
			url:  "https://i1.wp.com/dynamic-assets.imageg.top/uploads/798000/798727.gif?quality=99&h=250",
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elements, err := ParseURL(tt.url)
			if err != nil {
				t.Error(err)
			}

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
		SteamsLen int
	}
	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Image",
			url:  "https://htvassets.club/uploads/776000/776272.png",
			want: want{
				Title:     "776272",
				Type:      "image",
				SteamsLen: 1,
			},
		}, {
			name: "Gif",
			url:  "https://i1.wp.com/dynamic-assets.imageg.top/uploads/798000/798727.gif?quality=99&h=250",
			want: want{
				Title:     "798727",
				Type:      "gif",
				SteamsLen: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Extract(tt.url)
			if err != nil {
				t.Error(err)
			}

			want := want{
				Title:     data[0].Title,
				Type:      data[0].Type,
				SteamsLen: len(data[0].Streams),
			}
			if want != tt.want {
				t.Errorf("Got: %v - want: %v", want, tt.want)
			}
		})
	}
}
