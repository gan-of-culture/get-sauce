package utils

import (
	"reflect"
	"testing"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/test"
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

func TestParseHLSMaster(t *testing.T) {
	tests := []struct {
		Name   string
		master string
		Want   []*static.Stream
	}{
		{
			Name: "rich m3u master",
			master: `
			#EXTM3U
			#EXT-X-VERSION:4
			#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="audio_aac",NAME="Japanese",LANGUAGE="ja",AUTOSELECT=YES,DEFAULT=YES,URI="audio/aac/ja/stream.m3u8"
			
			# Media Playlists
			#EXT-X-STREAM-INF:AVERAGE-BANDWIDTH=3203863,BANDWIDTH=4634114,AUDIO="audio_aac"
			media-1/stream.m3u8
			#EXT-X-STREAM-INF:AVERAGE-BANDWIDTH=2189669,RESOLUTION=1280x720,AUDIO="audio_aac"
			media-2/stream.m3u8
			#EXT-X-STREAM-INF:AVERAGE-BANDWIDTH=1172227,BANDWIDTH=2479383,CODECS="avc1.42C01F,mp4a.40.2",RESOLUTION=864x486,AUDIO="audio_aac"
			media-3/stream.m3u8
			
			# I-Frame Playlists
			#EXT-X-I-FRAME-STREAM-INF:AVERAGE-BANDWIDTH=161759,BANDWIDTH=439418,CODECS="avc1.640032",RESOLUTION=1920x1080,URI="media-1/iframes.m3u8"
			#EXT-X-I-FRAME-STREAM-INF:AVERAGE-BANDWIDTH=120390,BANDWIDTH=262949,CODECS="avc1.4D401F",RESOLUTION=1280x720,URI="media-2/iframes.m3u8"
			#EXT-X-I-FRAME-STREAM-INF:AVERAGE-BANDWIDTH=53881,BANDWIDTH=146466,CODECS="avc1.42C01F",RESOLUTION=864x486,URI="media-3/iframes.m3u8"		
			`,
			Want: []*static.Stream{
				{
					Type: static.DataTypeVideo,
					URLs: []*static.URL{
						{
							URL: "media-1/stream.m3u8",
						},
					},
				}, {
					Type: static.DataTypeVideo,
					URLs: []*static.URL{
						{
							URL: "media-2/stream.m3u8",
						},
					},
					Quality: "1280x720",
				}, {
					Type: static.DataTypeVideo,
					URLs: []*static.URL{
						{
							URL: "media-3/stream.m3u8",
						},
					},
					Quality: "864x486",
					Info:    "avc1.42C01F,mp4a.40.2",
				}, {
					Type: static.DataTypeAudio,
					URLs: []*static.URL{
						{
							URL: "audio/aac/ja/stream.m3u8",
						},
					},
					Info: "ja",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			streams, err := ParseHLSMaster(&tt.master)
			test.CheckError(t, err)

			for i := range streams {
				if streams[i].Size != tt.Want[i].Size {
					t.Errorf("Got: %v - Want: %v", streams[i].Size, tt.Want[i].Size)
				}
				if streams[i].Quality != tt.Want[i].Quality {
					t.Errorf("Got: %v - Want: %v", streams[i].Size, tt.Want[i].Size)
				}
				if streams[i].Info != tt.Want[i].Info {
					t.Errorf("Got: %v - Want: %v", streams[i].Info, tt.Want[i].Info)
				}
				if streams[i].URLs[0].URL != tt.Want[i].URLs[0].URL {
					t.Errorf("Got: %v - Want: %v", streams[i].URLs[0].URL, tt.Want[i].URLs[0].URL)
				}
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
		in   string
		Want string
	}{
		{
			in:   "https://longurl.demo?fileext=mp4",
			Want: "mp4",
		}, {
			in:   "filename.mkv",
			Want: "mkv",
		},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			ext := GetFileExt(tt.in)

			if ext != tt.Want {
				t.Errorf("Got: %v - Want: %v", ext, tt.Want)
			}
		})
	}
}
