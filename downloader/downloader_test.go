package downloader

import (
	"testing"

	"github.com/gan-of-culture/go-hentai-scraper/static"
)

func TestDownload(t *testing.T) {
	tests := []struct {
		name string
		data static.Data
		want error
	}{
		{
			name: "underhentai single episode",
			data: static.Data{
				Site:  "https://underhentai.net",
				Title: "kiss-hug episode 01",
				Type:  "video",
				Streams: map[string]static.Stream{
					"0": {
						URLs: []static.URL{
							{
								URL: "https://sukebei.nyaa.si/download/2895133.torrent",
								Ext: "mp4",
							},
						},
					},
				},
				Url: "https://www.underhentai.net/kiss-hug/",
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
				Title: "Magical_Sempai_(Series) Magician_Sempai skyfreedom",
				Type:  "image",
				Streams: map[string]static.Stream{
					"0": {
						URLs: []static.URL{
							{
								URL: "https://scarlet.paheal.net/_images/a73b8b0053fd525488b1dbfd1b5ac2ed/3427635%20-%20Magical_Sempai_%28Series%29%20Magician_Sempai%20skyfreedom.jpg",
								Ext: "jpg",
							},
						},
					},
				},
				Url: "https://rule34.paheal.net/post/view/3427635",
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Download(tt.data)
			if err != tt.want {
				t.Error(err)
			}
		})
	}
}
