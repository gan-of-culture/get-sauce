package utils

import (
	"reflect"
	"testing"

	"github.com/gan-of-culture/go-hentai-scraper/config"
)

func TestGetLastItem(t *testing.T) {
	tests := []struct {
		name string
		list []string
		want string
	}{
		{
			name: "String slice",
			list: []string{"1", "2", "3", "last item"},
			want: "last item",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := GetLastItemString(tt.list)

			if item != tt.want {
				t.Errorf("Got: %v - want: %v", item, tt.want)
			}
		})
	}
}

func TestCalcSizeInByte(t *testing.T) {
	tests := []struct {
		name   string
		number float64
		unit   string
		want   int64
	}{
		{
			name:   "Kilobytes to Bytes",
			number: 752,
			unit:   "KB",
			want:   752000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes := CalcSizeInByte(tt.number, tt.unit)

			if bytes != tt.want {
				t.Errorf("Got: %v - want: %v", bytes, tt.want)
			}
		})
	}
}

func TestNeedDownloadList(t *testing.T) {
	type args struct {
		len int
	}
	tests := []struct {
		name  string
		args  args
		want  []int
		pages string
	}{
		{
			name: "pages test",
			args: args{
				len: 3,
			},
			pages: "1, 3",
			want:  []int{1, 3},
		},
		{
			name: "from to item selection 1",
			args: args{
				len: 10,
			},
			pages: "1-3, 5, 7-8, 10",
			want:  []int{1, 2, 3, 5, 7, 8, 10},
		},
		{
			name: "from to item selection 2",
			args: args{
				len: 10,
			},
			pages: "1,2, 4 , 5, 7-8  , 10",
			want:  []int{1, 2, 4, 5, 7, 8, 10},
		},
		{
			name: "from to item selection 3",
			args: args{
				len: 10,
			},
			pages: "5-1, 2",
			want:  []int{2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.Pages = tt.pages
			if got := NeedDownloadList(tt.args.len); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NeedDownloadList() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*func TestIsInTest(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "Test",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := IsInTests()

			if !is {
				t.Errorf("Got: %v - want: %v", is, tt.want)
			}
		})
	}
}*/

func TestGetMediaType(t *testing.T) {
	tests := []struct {
		ext  string
		want string
	}{
		{
			ext:  "jpg",
			want: "image/jpg",
		}, {
			ext:  "jpeg",
			want: "image/jpeg",
		}, {
			ext:  "png",
			want: "image/png",
		}, {
			ext:  "gif",
			want: "image/gif",
		}, {
			ext:  "webp",
			want: "image/webp",
		}, {
			ext:  "webm",
			want: "video/webm",
		}, {
			ext:  "mp4",
			want: "video/mp4",
		}, {
			ext:  "mkv",
			want: "video/mkv",
		}, {
			ext:  "m4a",
			want: "video/m4a",
		}, {
			ext:  "txt",
			want: "unknown/txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			dtype := GetMediaType(tt.ext)

			if dtype != tt.want {
				t.Errorf("Got: %v - want: %v", dtype, tt.want)
			}
		})
	}
}

func TestGetH1(t *testing.T) {
	tests := []struct {
		htmlString string
		want       string
	}{
		{
			htmlString: `<h1 class="entry-title" itemprop="name">Overflow 8</h1>`,
			want:       "Overflow 8",
		},
		{
			htmlString: `<h1>Overflow 8</h1>`,
			want:       "Overflow 8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.htmlString, func(t *testing.T) {
			h1 := GetH1(tt.htmlString)

			if h1 != tt.want {
				t.Errorf("Got: %v - want: %v", h1, tt.want)
			}
		})
	}
}
