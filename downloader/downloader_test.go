package downloader

import (
	"testing"

	"github.com/gan-of-culture/go-hentai-scraper/config"
	"github.com/gan-of-culture/go-hentai-scraper/static"
)

func TestDownload(t *testing.T) {
	tests := []struct {
		name string
		data static.Data
		want error
	}{
		{
			name: "hentaistream.moe 4k episode concurWriter",
			data: static.Data{
				Site:  "https://hentaistream.moe/",
				Title: "Overflow 1",
				Type:  "video",
				Streams: map[string]static.Stream{
					"0": {
						URLs: []static.URL{
							{
								URL: "https://01cdn.hentaistream.moe/2021/02/Overflow/E01/av1.2160p.webm",
								Ext: "webm",
							},
						},
						Size: int64(96865295),
					},
				},
			},
		}, {
			name: "rule34.xxx single img",
			data: static.Data{
				Site:  "https://rule34.xxx",
				Title: "4470590",
				Type:  "image/jpg",
				Streams: map[string]static.Stream{
					"0": {
						URLs: []static.URL{
							{
								URL: "https://wimg.rule34.xxx//images/3942/089a5ea08c47a1e79df5cb58b334693f686709de.jpg?4470590",
								Ext: "jpg",
							},
						},
					},
				},
				Url: "https://rule34.xxx/index.php?page=post&s=view&id=4470590",
			},
			want: nil,
		}, {
			name: "danbooru single post",
			data: static.Data{
				Site:  "https://danbooru.donmai.us/",
				Title: " touhou konpaku youmu niwashi  yuyu ",
				Type:  "image",
				Streams: map[string]static.Stream{
					"0": {
						URLs: []static.URL{
							{
								URL: "https://danbooru.donmai.us/data/sample/sample-3b63c93d7477967d0537d1d86d208b02.jpg",
								Ext: "jpg",
							},
						},
					},
				},
				Url: "https://danbooru.donmai.us/posts/3749687",
			},
			want: nil,
		}, {
			name: "rule 34 single post image",
			data: static.Data{
				Site:  "https://rule34.paheal.net",
				Title: "The_Dark_Mangaka tagme",
				Type:  "image",
				Streams: map[string]static.Stream{
					"0": {
						URLs: []static.URL{
							{
								URL: "https://lotus.paheal.net/_images/886fc1aeb1e67fedbbc2f1a40431fdc8/4322498%20-%20The_Dark_Mangaka%20tagme.png",
								Ext: "png",
							},
						},
					},
				},
				Url: "https://rule34.paheal.net/post/view/4322498",
			},
			want: nil,
		}, {
			name: "nhentai single page",
			data: static.Data{
				Site:  "https://nhentai.net",
				Title: "(C97) [H@BREAK (Itose Ikuto)] Koe Dashicha Barechau kara! [English]",
				Type:  "image",
				Streams: map[string]static.Stream{
					"0": {
						URLs: []static.URL{
							{
								URL: "https://i.nhentai.net/galleries/1550711/14.jpg",
								Ext: "jpg",
							},
						},
					},
				},
			},
		}, /*{
			name: "m3u8 normal",
			data: static.Data{
				Site:  "https://hentaistream.xxx/",
				Title: "Hime-sama Love Life! Episode 3",
				Type:  "application/x-mpegurl",
				Streams: map[string]static.Stream{
					"0": {
						URLs: []static.URL{
							{
								URL: "https://cdn1.htstreaming.com/cdn/down/7216c29dee7815942188208fe13e4068/360p/360p.txt",
								Ext: "mp4",
							},
						},
					},
				},
				Url: "https://hentaistream.xxx/watch/hime-sama-love-life-episode-3_P9TlY9FAOGHM7nn.html",
			},
		},*/{
			name: "m3u8 with aes-128 key",
			data: static.Data{
				Site:  "https://hanime.tv/",
				Title: "Toilet no Hanako-san vs Kukkyou Taimashi 2",
				Type:  "application/x-mpegurl",
				Streams: map[string]static.Stream{
					"0": {
						URLs: []static.URL{
							{
								URL: "https://m3u8s.highwinds-cdn.com/api/v9/m3u8s/j8gzqzcz3ck8pph771qqyr95v1ZzA049s9nkclk0pphclgbgn9sq2.m3u8",
								Ext: "ts",
							},
						},
					},
				},
				Url: "https://hanime.tv/videos/hentai/toilet-no-hanako-san-vs-kukkyou-taimashi-2",
			},
		},
	}
	config.Workers = 5
	downloader := New("0", false)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := downloader.Download(tt.data)
			if err != tt.want {
				t.Error(err)
			}
		})
	}
}
