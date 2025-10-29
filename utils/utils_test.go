package utils

import (
	"reflect"
	"testing"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/static"
)

func TestGetLastItem(t *testing.T) {
	tests := []struct {
		Name string
		list []string
		Want string
	}{
		{
			Name: "String slice",
			list: []string{"1", "2", "3", "last item"},
			Want: "last item",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			item := GetLastItemString(tt.list)

			if item != tt.Want {
				t.Errorf("Got: %v - Want: %v", item, tt.Want)
			}
		})
	}
}

func TestCalcSizeInByte(t *testing.T) {
	tests := []struct {
		Name   string
		number float64
		unit   string
		Want   int64
	}{
		{
			Name:   "Kilobytes to Bytes",
			number: 752,
			unit:   "KB",
			Want:   752000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			bytes := CalcSizeInByte(tt.number, tt.unit)

			if bytes != tt.Want {
				t.Errorf("Got: %v - Want: %v", bytes, tt.Want)
			}
		})
	}
}

func TestByteCountSI(t *testing.T) {
	tests := []struct {
		Name   string
		number int64
		Want   string
	}{
		{
			Name:   "To Kilobytes",
			number: 752000,
			Want:   "752.0 kB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			bytes := ByteCountSI(tt.number)

			if bytes != tt.Want {
				t.Errorf("Got: %v - Want: %v", bytes, tt.Want)
			}
		})
	}
}

func TestNeedDownloadList(t *testing.T) {
	type args struct {
		len int
	}
	tests := []struct {
		Name  string
		args  args
		Want  []int
		pages string
	}{
		{
			Name: "pages test",
			args: args{
				len: 3,
			},
			pages: "1, 3",
			Want:  []int{1, 3},
		},
		{
			Name: "from to item selection 1",
			args: args{
				len: 10,
			},
			pages: "1-3, 5, 7-8, 10",
			Want:  []int{1, 2, 3, 5, 7, 8, 10},
		},
		{
			Name: "from to item selection 2",
			args: args{
				len: 10,
			},
			pages: "1,2, 4 , 5, 7-8  , 10",
			Want:  []int{1, 2, 4, 5, 7, 8, 10},
		},
		{
			Name: "from to item selection 3",
			args: args{
				len: 10,
			},
			pages: "5-1, 2",
			Want:  []int{2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			config.Pages = tt.pages
			if got := NeedDownloadList(tt.args.len); !reflect.DeepEqual(got, tt.Want) {
				t.Errorf("NeedDownloadList() = %v, Want %v", got, tt.Want)
			}
		})
	}
}

func TestGetMediaType(t *testing.T) {
	tests := []struct {
		ext  string
		Want static.DataType
	}{
		{
			ext:  "jpg",
			Want: "image",
		}, {
			ext:  "jpeg",
			Want: "image",
		}, {
			ext:  "png",
			Want: "image",
		}, {
			ext:  "gif",
			Want: "image",
		}, {
			ext:  "webp",
			Want: "image",
		}, {
			ext:  "webm",
			Want: "video",
		}, {
			ext:  "mp4",
			Want: "video",
		}, {
			ext:  "mkv",
			Want: "video",
		}, {
			ext:  "m4a",
			Want: "video",
		}, {
			ext:  "txt",
			Want: "video",
		}, {
			ext:  "m3u8",
			Want: "video",
		},
	}
	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			dtype := GetMediaType(tt.ext)

			if dtype != tt.Want {
				t.Errorf("Got: %v - Want: %v", dtype, tt.Want)
			}
		})
	}
}

func TestGetH1(t *testing.T) {
	tests := []struct {
		Name       string
		htmlString string
		idx        int
		Want       string
	}{
		{
			Name:       "h1 tag with params",
			htmlString: `<h1 class="entry-title" itemprop="name">Overflow 8</h1>`,
			idx:        0,
			Want:       "Overflow 8",
		},
		{
			Name:       "normal case",
			htmlString: `<h1>Overflow 8</h1>`,
			idx:        0,
			Want:       "Overflow 8",
		},
		{
			Name:       "get specific",
			htmlString: `<h1>Overflow 8</h1><h1>Overflow 9</h1><h1>Overflow 10</h1>`,
			idx:        1,
			Want:       "Overflow 9",
		}, {
			Name:       "out of range but has one, return it",
			htmlString: `<h1>Overflow 8</h1>`,
			idx:        1,
			Want:       "Overflow 8",
		}, {
			Name:       "last",
			htmlString: `<h1>Overflow 7</h1><h1>Overflow 8</h1>`,
			idx:        -1,
			Want:       "Overflow 8",
		}, {
			Name:       "escaped character",
			htmlString: `<h1 class="title">Queen&#39;s Discipline</h1>`,
			idx:        -1,
			Want:       "Queen's Discipline",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			h1 := GetH1(&tt.htmlString, tt.idx)

			if h1 != tt.Want {
				t.Errorf("Got: %v - Want: %v", h1, tt.Want)
			}
		})
	}
}

func TestGetSectionHeadingElement(t *testing.T) {
	tests := []struct {
		Name       string
		htmlString string
		level      int
		idx        int
		Want       string
	}{
		{
			Name:       "h1",
			htmlString: `<h1>Overflow 8</h1>`,
			level:      1,
			idx:        0,
			Want:       "Overflow 8",
		},
		{
			Name:       "h2",
			htmlString: `<h2>Overflow 8</h2>`,
			level:      2,
			idx:        0,
			Want:       "Overflow 8",
		},
		{
			Name:       "h3",
			htmlString: `<h3>Overflow 8</h3>`,
			level:      3,
			idx:        0,
			Want:       "Overflow 8",
		},
		{
			Name:       "h4",
			htmlString: `<h4>Overflow 8</h4>`,
			level:      4,
			idx:        0,
			Want:       "Overflow 8",
		},
		{
			Name:       "h5",
			htmlString: `<h5>Overflow 8</h5>`,
			level:      5,
			idx:        0,
			Want:       "Overflow 8",
		},
		{
			Name:       "h6",
			htmlString: `<h6>Overflow 8</h6>`,
			level:      6,
			idx:        0,
			Want:       "Overflow 8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			h1 := GetSectionHeadingElement(&tt.htmlString, tt.level, tt.idx)

			if h1 != tt.Want {
				t.Errorf("Got: %v - Want: %v", h1, tt.Want)
			}
		})
	}
}

func TestMeta(t *testing.T) {
	tests := []struct {
		htmlString string
		property   string
		Want       string
	}{
		{
			htmlString: `<meta property="og:title" content="Imouto Paradise! 3 The Animation Episode 1" />`,
			property:   "og:title",
			Want:       "Imouto Paradise! 3 The Animation Episode 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.htmlString, func(t *testing.T) {
			h1 := GetMeta(&tt.htmlString, tt.property)

			if h1 != tt.Want {
				t.Errorf("Got: %v - Want: %v", h1, tt.Want)
			}
		})
	}
}

func TestRemoveAdjDuplicates(t *testing.T) {
	tests := []struct {
		Name string
		in   []string
		Want []string
	}{
		{
			Name: "default",
			in:   []string{"test", "test", "hello", "world"},
			Want: []string{"test", "hello", "world"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			s := RemoveAdjDuplicates(tt.in)

			if len(s) != len(tt.Want) {
				t.Errorf("Got: %v - Want: %v", tt.in, tt.Want)
			}
		})
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		Name string
		pair struct {
			err error
			ctx string
		}
		Want string
	}{
		{
			Name: "String slice",
			pair: struct {
				err error
				ctx string
			}{
				static.ErrURLParseFailed,
				"https://google.com",
			},
			Want: static.ErrURLParseFailed.Error() + ": " + "https://google.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := Wrap(tt.pair.err, tt.pair.ctx)

			if err.Error() != tt.Want {
				t.Errorf("Got: %v - Want: %v", err.Error(), tt.Want)
			}
		})
	}
}

func TestGetFileExt(t *testing.T) {
	tests := []struct {
		Name string
		In   string
		Want string
	}{
		{
			Name: "From URL query",
			In:   "https://longurl.demo?fileext=mp4",
			Want: "mp4",
		}, {
			Name: "Default",
			In:   "filename.mkv",
			Want: "mkv",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			ext := GetFileExt(tt.In)

			if ext != tt.Want {
				t.Errorf("Got: %v - Want: %v", ext, tt.Want)
			}
		})
	}
}

func TestSortStreamsBySize(t *testing.T) {
	tests := []struct {
		Name string
		In   map[string]*static.Stream
		Want map[string]*static.Stream
	}{
		{
			Name: "Unsorted",
			In: map[string]*static.Stream{
				"0": {
					Size: 1234567,
				},
				"1": {
					Size: 123456,
				},
				"2": {
					Size: 12345678,
				},
			},
			Want: map[string]*static.Stream{
				"0": {
					Size: 12345678,
				},
				"1": {
					Size: 1234567,
				},
				"2": {
					Size: 123456,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			streams := SortStreamsBySize(tt.In)

			if !reflect.DeepEqual(streams, tt.Want) {
				t.Errorf("Got: %v - Want: %v", streams, tt.Want)
			}
		})
	}
}
