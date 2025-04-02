package downloader

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/config"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/test"
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
								URL: "https://r34i.paheal-cdn.net/cf/21/cf21c36b64db166b1e1aac9f3243d3ec",
								Ext: "jpg",
							},
						},
					},
				},
				URL: "https://rule34.paheal.net/post/view/4698365",
			},
			Want: nil,
		},
	}
	config.Workers = 5
	config.SelectStream = "0"
	config.Keep = true
	config.Truncate = true
	downloader := New(false)
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			err := downloader.Download(tt.data)
			test.CheckError(t, err)
		})
	}
}
