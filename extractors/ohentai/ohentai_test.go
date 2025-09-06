package ohentai

// func TestParseURL(t *testing.T) {
// 	tests := []struct {
// 		Name string
// 		URL  string
// 		Want int
// 	}{
// 		{
// 			Name: "Single Video",
// 			URL:  "https://ohentai.org/detail.php?vid=MzM5MQ==",
// 			Want: 1,
// 		}, {
// 			Name: "Tag",
// 			URL:  "https://ohentai.org/tagsearch.php?tag=Uncensored",
// 			Want: 24,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.Name, func(t *testing.T) {
// 			URLs := parseURL(tt.URL)
// 			if len(URLs) > tt.Want || len(URLs) == 0 {
// 				t.Errorf("Got: %v - Want: %v", len(URLs), tt.Want)
// 			}
// 		})
// 	}
// }

// func TestExtract(t *testing.T) {
// 	tests := []struct {
// 		Name string
// 		Args test.Args
// 	}{
// 		{
// 			Name: "Single Video",
// 			Args: test.Args{
// 				URL:     "https://ohentai.org/detail.php?vid=MzM5MQ==",
// 				Title:   "Shirakami Collection - Episode 2",
// 				Quality: "",
// 				Size:    143959503,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.Name, func(t *testing.T) {
// 			data, err := New().Extract(tt.Args.URL)
// 			test.CheckError(t, err)
// 			test.Check(t, tt.Args, data[0])
// 		})
// 	}
// }
