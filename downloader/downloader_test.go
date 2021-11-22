package downloader

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/static"
)

func TestDownload(t *testing.T) {
	tests := []struct {
		Name string
		data *static.Data
		Want error
	}{
		{
			Name: "hentaistream.moe 4k episode concurWriter",
			data: &static.Data{
				Site:  "https://hentaistream.moe/",
				Title: "Overflow 1",
				Type:  "video",
				Streams: map[string]*static.Stream{
					"0": {
						URLs: []*static.URL{
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
			Name: "rule34.xxx single img",
			data: &static.Data{
				Site:  "https://rule34.xxx",
				Title: "4470590",
				Type:  "image/jpg",
				Streams: map[string]*static.Stream{
					"0": {
						URLs: []*static.URL{
							{
								URL: "https://wimg.rule34.xxx//images/3942/089a5ea08c47a1e79df5cb58b334693f686709de.jpg?4470590",
								Ext: "jpg",
							},
						},
					},
				},
				URL: "https://rule34.xxx/index.php?page=post&s=view&id=4470590",
			},
			Want: nil,
		}, {
			Name: "danbooru single post",
			data: &static.Data{
				Site:  "https://danbooru.donmai.us/",
				Title: " touhou konpaku youmu niwashi  yuyu ",
				Type:  "image",
				Streams: map[string]*static.Stream{
					"0": {
						URLs: []*static.URL{
							{
								URL: "https://danbooru.donmai.us/data/sample/sample-3b63c93d7477967d0537d1d86d208b02.jpg",
								Ext: "jpg",
							},
						},
					},
				},
				URL: "https://danbooru.donmai.us/posts/3749687",
			},
			Want: nil,
		}, {
			Name: "rule 34 single post image",
			data: &static.Data{
				Site:  "https://rule34.paheal.net",
				Title: "Ahri Cian_Yo League_of_Legends",
				Type:  "image",
				Streams: map[string]*static.Stream{
					"0": {
						URLs: []*static.URL{
							{
								URL: "https://peach.paheal.net/_images/cf21c36b64db166b1e1aac9f3243d3ec/4698365%20-%20Ahri%20Cian_Yo%20League_of_Legends.jpg",
								Ext: "jpg",
							},
						},
					},
				},
				URL: "https://rule34.paheal.net/post/view/4698365",
			},
			Want: nil,
		}, {
			Name: "nhentai single page",
			data: &static.Data{
				Site:  "https://nhentai.net",
				Title: "(C97) [H@BREAK (Itose Ikuto)] Koe Dashicha Barechau kara! [English]",
				Type:  "image",
				Streams: map[string]*static.Stream{
					"0": {
						URLs: []*static.URL{
							{
								URL: "https://i.nhentai.net/galleries/1550711/14.jpg",
								Ext: "jpg",
							},
						},
					},
				},
			},
		}, /*{
			Name: "m3u8 normal",
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
			Name: "m3u8 with aes-128 key - this is not a complete file",
			data: &static.Data{
				Site:  "https://hanime.tv/",
				Title: "Toilet no Hanako-san vs Kukkyou Taimashi 2",
				Type:  "video",
				Streams: map[string]*static.Stream{
					"0": {
						URLs: []*static.URL{
							{
								URL: "https://new.alphafish.top/2/8/0/9/v1x/segs/b0/0/G6Q0WPempZmwKwgPIPEwD3hW.html",
								Ext: "ts",
							},
							{
								URL: "https://order.apperoni.top/2/8/0/9/v1x/segs/b0/0/3PsK4HaXBX1MEcTE6pxNbo3T.html",
								Ext: "ts",
							},
							{
								URL: "https://dash.blingo.top/2/8/0/9/v1x/segs/b0/0/kFV4nuFjh9wTjlzImekNqShk.html",
								Ext: "ts",
							},
							{
								URL: "https://portal.nodebook11.top/2/8/0/9/v1x/segs/b0/0/8B9WLdrUH0oGWIcPpwN7BNRt.html",
								Ext: "ts",
							},
							{
								URL: "https://new.alphafish.top/2/8/0/9/v1x/segs/b0/0/93X7sk6FSwhtSRpZnCF5vodk.html",
								Ext: "ts",
							},
						},
						Ext: "ts",
						Key: []byte{48, 49, 50, 51, 52, 53, 54, 55, 48, 49, 50, 51, 52, 53, 54, 55},
					},
				},
				URL: "https://hanime.tv/videos/hentai/toilet-no-hanako-san-vs-kukkyou-taimashi-2",
			},
		},
	}
	config.Workers = 5
	config.SelectStream = "0"
	downloader := New(false)
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := downloader.Download(tt.data)
			if err != tt.Want {
				t.Error(err)
			}
		})
	}
}
