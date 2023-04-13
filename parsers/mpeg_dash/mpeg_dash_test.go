package mpegdash

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/test"
)

func TestReplaceIdentifier(t *testing.T) {
	tests := []struct {
		Name string
		URI  string
		Vars map[string]string
		Want string
	}{
		{
			Name: "No identifier",
			URI:  "init-stream.html",
			Vars: map[string]string{},
			Want: "init-stream.html",
		},
		{
			Name: "Single",
			URI:  `chunks\init-stream$RepresentationID$.html`,
			Vars: map[string]string{
				"RepresentationID": "0",
			},
			Want: `chunks\init-stream0.html`,
		},
		{
			Name: "With format tag",
			URI:  `chunks\chunk-stream$RepresentationID$-$Number%05d$.html`,
			Vars: map[string]string{
				"RepresentationID": "0",
				"Number":           "1",
			},
			Want: `chunks\chunk-stream0-00001.html`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			uri, err := replaceIdentifier(tt.URI, tt.Vars)
			if uri != tt.Want {
				t.Errorf("Got: %v - Want: %v", uri, tt.Want)
			}
			test.CheckError(t, err)
		})
	}
}

func TestExtractDASHManifest(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Args static.Stream
	}{
		{
			Name: "Default",
			URL:  "https://str.h-dl.xyz/2023/Class.de.Otoko.wa.Boku.Hitori/E01/2160/manifest.mpd",
			Args: static.Stream{
				Type:    static.DataTypeVideo,
				Quality: "3840x2160",
				Ext:     "mp4",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			streams, err := ExtractDASHManifest(tt.URL, nil)
			stream := streams["0"]
			test.CheckError(t, err)
			if stream.Type != tt.Args.Type {
				t.Errorf("Got: %v - Want: %v", stream.Type, tt.Args.Type)
			}
			if stream.Ext != tt.Args.Ext {
				t.Errorf("Got: %v - Want: %v", stream.Ext, tt.Args.Ext)
			}
			if stream.Quality != tt.Args.Quality {
				t.Errorf("Got: %v - Want: %v", stream.Quality, tt.Args.Quality)
			}
		})
	}
}
