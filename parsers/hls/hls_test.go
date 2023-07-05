package hls

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/test"
)

func TestParseMaster(t *testing.T) {
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
			streams, err := ParseMaster(&tt.master)
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

func TestExtract(t *testing.T) {
	tests := []struct {
		Name    string
		URL     string
		Headers map[string]string
		Want    map[string]*static.Stream
	}{
		{
			Name: "HLS where stream order is from small to high",
			URL:  "https://na-02.javprovider.com/hls/K/kuroinu-ii-animation/1/playlist.m3u8",
			Headers: map[string]string{
				"Referer": "https://hentaimama.io",
			},
			Want: map[string]*static.Stream{
				"0": {
					Type:    static.DataTypeVideo,
					Quality: "1280x720",
				},
				"1": {
					Type:    static.DataTypeVideo,
					Quality: "842x480",
				},
				"2": {
					Type:    static.DataTypeVideo,
					Quality: "640x360",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			streams, err := Extract(tt.URL, tt.Headers)
			test.CheckError(t, err)
			for k, v := range streams {
				if v.Quality != tt.Want[k].Quality {
					t.Errorf("Got: %v - Want: %v", v.Quality, tt.Want[k].Quality)
				}
			}
		})
	}
}
