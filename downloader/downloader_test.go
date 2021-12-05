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
			Name: "muchohentai.com download video + audio + captions separately",
			data: &static.Data{
				Site:  "https://muchohentai.com/",
				Title: "Onaho Kyoushitsu: Joshi Zenin Ninshin Keikaku The Animation Episode 1 English Subbed",
				Type:  static.DataTypeVideo,
				Streams: map[string]*static.Stream{
					"0": {
						Type: static.DataTypeVideo,
						URLs: []*static.URL{
							{
								URL: "https://va03-edge.tmncdn.io/wp-content/uploads/Onaho_Kyoushitsu/episode_1/media-2/segment-0.ts",
								Ext: "ts",
							},
						},
						Quality: "1920x1080",
						Ext:     "mp4",
						Key:     []byte{112, 176, 72, 44, 26, 39, 59, 183, 219, 153, 186, 209, 70, 90, 13, 160},
					},
					"1": {
						Type: static.DataTypeAudio,
						URLs: []*static.URL{
							{
								URL: "https://va03-edge.tmncdn.io/wp-content/uploads/Onaho_Kyoushitsu/episode_1/audio/aac/ja/segment-0.aac",
								Ext: "aac",
							},
						},
						Ext: "aac",
						Key: []byte{112, 176, 72, 44, 26, 39, 59, 183, 219, 153, 186, 209, 70, 90, 13, 160},
					},
				},
				Captions: []*static.Caption{
					{
						URL: static.URL{
							URL: "https://muchohentai.com/wp-content/uploads/Onaho_Kyoushitsu/episode_1/subs/en.vtt",
							Ext: "vtt",
						},
						Language: "English",
					},
				},
				URL: "https://muchohentai.com/aBo4Rk/167062/",
			},
		}, {
			Name: "hentaistream.moe 4k episode concurWriter",
			data: &static.Data{
				Site:  "https://hentaistream.moe/",
				Title: "Overflow 1",
				Type:  static.DataTypeVideo,
				Streams: map[string]*static.Stream{
					"0": {
						Type: static.DataTypeVideo,
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
				Type:  static.DataTypeImage,
				Streams: map[string]*static.Stream{
					"0": {
						Type: static.DataTypeImage,
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
				Type:  static.DataTypeImage,
				Streams: map[string]*static.Stream{
					"0": {
						Type: static.DataTypeImage,
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
				Type:  static.DataTypeImage,
				Streams: map[string]*static.Stream{
					"0": {
						Type: static.DataTypeImage,
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
				Type:  static.DataTypeImage,
				Streams: map[string]*static.Stream{
					"0": {
						Type: static.DataTypeImage,
						URLs: []*static.URL{
							{
								URL: "https://i.nhentai.net/galleries/1550711/14.jpg",
								Ext: "jpg",
							},
						},
					},
				},
			},
		}, {
			Name: "m3u8 with aes-128 key - this is not a complete file",
			data: &static.Data{
				Site:  "https://hanime.tv/",
				Title: "Toilet no Hanako-san vs Kukkyou Taimashi 2",
				Type:  static.DataTypeVideo,
				Streams: map[string]*static.Stream{
					"0": {
						Type: static.DataTypeVideo,
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
	config.NoMerge = true
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
