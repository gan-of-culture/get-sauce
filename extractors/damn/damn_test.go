package damn

/*import (
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
			Name: "Single Episode",
			URL:  "https://www.damn.stream/watch/elf-hime-nina-episode-3",
			Want: 1,
		}, {
			Name: "Series",
			URL:  "https://www.damn.stream/hentai/ane-koi-suki-kirai-daisuki",
			Want: 2,
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
			Name: "Single Episode",
			Args: test.Args{
				URL:     "https://www.damn.stream/watch/elf-hime-nina-episode-3",
				Title:   "Corrupted End - Elf-hime Nina Episode 3",
				Quality: "",
				Size:    49962683,
			},
		}, {
			Name: "Series",
			Args: test.Args{
				URL:     "https://www.damn.stream/hentai/ane-koi-suki-kirai-daisuki",
				Title:   "Ane Koi: Suki Kirai Daisuki. Episode 2",
				Quality: "",
				Size:    58222658,
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
}*/
