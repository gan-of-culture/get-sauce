package hentaipulse

/*
func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single Episode",
			URL:  "https://hentaipulse.com/bonyuu-chan-wa-dashitai-episode-02-english-subbed/",
			Want: 1,
		}, {
			Name: "Overview",
			URL:  "https://hentaipulse.com/hentai-anime/english-subbed-hentai-anime/",
			Want: 24,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			URLs := parseURL(tt.URL)
			if len(URLs) != tt.Want || len(URLs) == 0 {
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
				URL:     "https://hentaipulse.com/bonyuu-chan-wa-dashitai-episode-02-english-subbed/",
				Title:   "bonyuu-chan-wa-dashitai-episode-02-english-subbed",
				Quality: "",
				Size:    96398241,
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
