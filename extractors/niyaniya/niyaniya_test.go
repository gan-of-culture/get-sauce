package niyaniya

/* func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want int
	}{
		{
			Name: "Single",
			URL:  "https://niyaniya.moe/g/24687/4b7e8e4b1936",
			Want: 1,
		},
		{
			Name: "Overview",
			URL:  "https://niyaniya.moe/tag/artist:alp",
			Want: 28,
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
			Name: "Single extraction",
			Args: test.Args{
				URL:     "https://niyaniya.moe/g/13848/85a0f534cb44",
				Title:   "[Alp] Reward Poolside (Comic Bavel 2016-08)",
				Quality: "1360x1920",
				Size:    17848085,
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
} */
