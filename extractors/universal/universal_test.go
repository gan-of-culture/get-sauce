package universal

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
			Name: "Imigur",
			Args: test.Args{
				URL:     "https://i.imgur.com/CVTp6cK.jpeg",
				Title:   "CVTp6cK",
				Quality: "",
				Size:    384586,
			},
		},
		{
			Name: "with bloat after ext",
			Args: test.Args{
				URL:     "https://wimg.rule34.xxx//images/41/d54d1f1dde83ff670c5932ab8f2d42a5cf18e587.jpg?40602",
				Title:   "d54d1f1dde83ff670c5932ab8f2d42a5cf18e587",
				Quality: "",
				Size:    125079,
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
