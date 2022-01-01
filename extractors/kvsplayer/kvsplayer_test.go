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
				URL:     "https://www.kvs-demo.com/videos/105/kelis-4th-of-july/",
				Title:   "kelis-4th-of-july",
				Quality: "720p",
				Size:    104777803,
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
