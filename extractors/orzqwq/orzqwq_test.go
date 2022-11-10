package orzqwq

import (
	"testing"

	"github.com/gan-of-culture/get-sauce/test"
)

func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single Gallery",
			URL:  "https://orzqwq.com/manga/pixiv-%e3%81%bf%e3%82%8c%e3%81%84%ef%bc%a0%f0%9f%98%88-746841/",
			Want: 1,
		}, {
			Name: "Tag",
			URL:  "https://orzqwq.com/manga-tag/original/",
			Want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			URLs := parseURL(tt.URL)
			if len(URLs) > tt.Want || len(URLs) == 0 {
				t.Errorf("Got: %v - Want: %v", len(URLs), tt.Want)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		Name string
		Args test.Args
	}{
		{
			Name: "Single Gallery",
			Args: test.Args{
				URL:     "https://orzqwq.com/manga/%E5%96%B5%E5%B0%8F%E5%90%8910%E6%9C%88-%E9%95%9C%E8%8A%B1%E6%B0%B4%E6%9C%88/",
				Title:   "喵小吉10月 镜花水月",
				Quality: "",
				Size:    0,
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
