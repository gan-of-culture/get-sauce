package jwplayer

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
			Name: "Single Episode hentai.guru",
			Args: test.Args{
				URL:     "https://hentai.guru/wp-content/plugins/player-logic/player.php?data=bDFCblBOMTZvcVAwZ2E2M1VMeUpldUc2MzNpeGUwZEc5dGF5RlRVWlBnbjJXQWV2YzN4RmJwemJjVno1TUYvWXRvMnNaNjhPMEdDdDJ1RlA0ZFVZYyt1WksybnFsd2lxZW1pQzJMSzYzd05Vdk1FVzEyeStaS3c1ekpqSnFaNmlPZE96UUM2VzljdlRJb0Zkc0tqVlJnPT06fDo6fDpzN0t3UmFNcnRNYUdJa3FrbENKSzhnPT0=",
				Title:   "jwplayer video",
				Quality: "1920x1080",
				Size:    475164360,
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

func TestFindJWPlayerURL(t *testing.T) {
	tests := []struct {
		Name string
		URL  string
		Want string
	}{
		{
			Name: "HTML string",
			URL:  `<iframe src="https://hentaihaven.xxx/wp-content/plugins/player-logic/player.php?data=OTZPUXF5OU1DbWxVSzROajdkWkx6T1Y4RFBVZ3JQKzRqWlB3M1FXY2tudUU2S1VHNklsUHY0TmRWREtCd3pKb2Q3ZVVtNklCWDJRUndqNGZQWjR5MjJKOEo2ZXU3QzBmTEU5aTRKK2tGME1MazZZUzh4Y0ZMRThyOG9TOTlHMnM0TnZvQnRIelhHTTJrL1htYUVjUEFRPT06fDo6fDp2bFlyazVWNmdCaU1BMTVOM0U4SGNBPT0=" frameborder="0" scrolling="no" allowfullscreen=""></iframe>`,
			Want: "https://hentaihaven.xxx/wp-content/plugins/player-logic/player.php?data=OTZPUXF5OU1DbWxVSzROajdkWkx6T1Y4RFBVZ3JQKzRqWlB3M1FXY2tudUU2S1VHNklsUHY0TmRWREtCd3pKb2Q3ZVVtNklCWDJRUndqNGZQWjR5MjJKOEo2ZXU3QzBmTEU5aTRKK2tGME1MazZZUzh4Y0ZMRThyOG9TOTlHMnM0TnZvQnRIelhHTTJrL1htYUVjUEFRPT06fDo6fDp2bFlyazVWNmdCaU1BMTVOM0U4SGNBPT0=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			u := FindJWPlayerURL(&tt.URL)
			if u == "" {
				t.Errorf("Got: %v - Want: %v", u, tt.Want)
			}
		})
	}
}
