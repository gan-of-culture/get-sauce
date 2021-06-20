package utils

import (
	"reflect"
	"testing"

	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/static"
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

func TestGetMediaType(t *testing.T) {
	tests := []struct {
		ext  string
		want static.DataType
	}{
		{
			ext:  "jpg",
			want: "image",
		}, {
			ext:  "jpeg",
			want: "image",
		}, {
			ext:  "png",
			want: "image",
		}, {
			ext:  "gif",
			want: "image",
		}, {
			ext:  "webp",
			want: "image",
		}, {
			ext:  "webm",
			want: "video",
		}, {
			ext:  "mp4",
			want: "video",
		}, {
			ext:  "mkv",
			want: "video",
		}, {
			ext:  "m4a",
			want: "video",
		}, {
			ext:  "txt",
			want: "video",
		}, {
			ext:  "m3u8",
			want: "video",
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
		name       string
		htmlString string
		idx        int
		want       string
	}{
		{
			name:       "h1 tag with params",
			htmlString: `<h1 class="entry-title" itemprop="name">Overflow 8</h1>`,
			idx:        0,
			want:       "Overflow 8",
		},
		{
			name:       "normal case",
			htmlString: `<h1>Overflow 8</h1>`,
			idx:        0,
			want:       "Overflow 8",
		},
		{
			name:       "get specific",
			htmlString: `<h1>Overflow 8</h1><h1>Overflow 9</h1><h1>Overflow 10</h1>`,
			idx:        1,
			want:       "Overflow 9",
		}, {
			name:       "out of range but has one, return it",
			htmlString: `<h1>Overflow 8</h1>`,
			idx:        1,
			want:       "Overflow 8",
		}, {
			name:       "last",
			htmlString: `<h1>Overflow 7</h1><h1>Overflow 8</h1>`,
			idx:        -1,
			want:       "Overflow 8",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h1 := GetH1(&tt.htmlString, tt.idx)

			if h1 != tt.want {
				t.Errorf("Got: %v - want: %v", h1, tt.want)
			}
		})
	}
}

func TestMeta(t *testing.T) {
	tests := []struct {
		htmlString string
		property   string
		want       string
	}{
		{
			htmlString: `<meta property="og:title" content="Imouto Paradise! 3 The Animation Episode 1" />`,
			property:   "og:title",
			want:       "Imouto Paradise! 3 The Animation Episode 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.htmlString, func(t *testing.T) {
			h1 := GetMeta(&tt.htmlString, tt.property)

			if h1 != tt.want {
				t.Errorf("Got: %v - want: %v", h1, tt.want)
			}
		})
	}
}

func TestRemoveAdjDuplicates(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want []string
	}{
		{
			name: "default",
			in:   []string{"test", "test", "hello", "world"},
			want: []string{"test", "hello", "world"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := RemoveAdjDuplicates(tt.in)

			if len(s) != len(tt.want) {
				t.Errorf("Got: %v - want: %v", tt.in, tt.want)
			}
		})
	}
}

func TestParseM3UMaster(t *testing.T) {
	tests := []struct {
		name   string
		master string
		want   map[string]*static.Stream
	}{
		{
			name: "rich m3u master",
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
			want: map[string]*static.Stream{
				"0": {
					URLs: []*static.URL{
						{
							URL: "media-1/stream.m3u8",
						},
					},
					Size: 4634114,
				},
				"1": {
					URLs: []*static.URL{
						{
							URL: "media-2/stream.m3u8",
						},
					},
					Quality: "1280x720",
				},
				"2": {
					URLs: []*static.URL{
						{
							URL: "media-3/stream.m3u8",
						},
					},
					Size:    2479383,
					Quality: "864x486",
					Info:    "avc1.42C01F,mp4a.40.2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			streams, err := ParseM3UMaster(&tt.master)
			if err != nil {
				t.Error(err)
			}

			for i := range streams {
				if streams[i].Size != tt.want[i].Size {
					t.Errorf("Got: %v - want: %v", streams[i].Size, tt.want[i].Size)
				}
				if streams[i].Quality != tt.want[i].Quality {
					t.Errorf("Got: %v - want: %v", streams[i].Size, tt.want[i].Size)
				}
				if streams[i].Info != tt.want[i].Info {
					t.Errorf("Got: %v - want: %v", streams[i].Info, tt.want[i].Info)
				}
				if streams[i].URLs[0].URL != tt.want[i].URLs[0].URL {
					t.Errorf("Got: %v - want: %v", streams[i].URLs[0].URL, tt.want[i].URLs[0].URL)
				}
			}
		})
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		name string
		pair struct {
			err error
			ctx string
		}
		want string
	}{
		{
			name: "String slice",
			pair: struct {
				err error
				ctx string
			}{
				static.ErrURLParseFailed,
				"https://google.com",
			},
			want: static.ErrURLParseFailed.Error() + ": " + "https://google.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Wrap(tt.pair.err, tt.pair.ctx)

			if err.Error() != tt.want {
				t.Errorf("Got: %v - want: %v", err.Error(), tt.want)
			}
		})
	}
}
