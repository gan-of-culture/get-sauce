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
			Name: "awwni",
			Args: test.Args{
				URL:     "http://cdn.awwni.me/16c8v.jpg",
				Title:   "16c8v",
				Quality: "",
				Size:    102169,
			},
		},
		{
			Name: "with bloat after ext",
			Args: test.Args{
				URL:     "https://img.rule34.xxx//images/1979/b84be533024a3d1dcc6b01c0cb7358c9.jpeg?2686173",
				Title:   "b84be533024a3d1dcc6b01c0cb7358c9",
				Quality: "",
				Size:    143583,
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
