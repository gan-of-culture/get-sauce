package kvsplayer

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/test"
)

func TestExtract(t *testing.T) {
	tests := []struct {
		Name string
		Args test.Args
	}{
		{
			Name: "Single Episode www.kvs-demo.com",
			Args: test.Args{
				URL:     "https://hentaibar.com/videos/2688/soshite-watashi-wa-ojisan-ni-episode-4-english-subbed/",
				Title:   "soshite-watashi-wa-ojisan-ni-episode-4-english-subbed",
				Quality: "1080p",
				Size:    550603331,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.Args.URL)
			test.CheckError(t, err)
			test.Check(t, tt.Args, data[0])
		})
	}
}
