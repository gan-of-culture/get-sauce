package jwplayer

import "testing"

func TestExtract(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Episode hentai.guru",
			url:  "https://hentai.guru/wp-content/plugins/player-logic/player.php?data=bDFCblBOMTZvcVAwZ2E2M1VMeUpldUc2MzNpeGUwZEc5dGF5RlRVWlBnbjJXQWV2YzN4RmJwemJjVno1TUYvWXRvMnNaNjhPMEdDdDJ1RlA0ZFVZYyt1WksybnFsd2lxZW1pQzJMSzYzd05Vdk1FVzEyeStaS3c1ekpqSnFaNmlPZE96UUM2VzljdlRJb0Zkc0tqVlJnPT06fDo6fDpzN0t3UmFNcnRNYUdJa3FrbENKSzhnPT0=",
			want: 1,
		},
		/*{
			name: "Single Episode hentaihaven.xxx",
			url:  "https://hentaihaven.xxx/wp-content/plugins/player-logic/player.php?data=OTZPUXF5OU1DbWxVSzROajdkWkx6T1Y4RFBVZ3JQKzRqWlB3M1FXY2tudUU2S1VHNklsUHY0TmRWREtCd3pKb2Q3ZVVtNklCWDJRUndqNGZQWjR5MjJKOEo2ZXU3QzBmTEU5aTRKK2tGME1MazZZUzh4Y0ZMRThyOG9TOTlHMnM0TnZvQnRIelhHTTJrL1htYUVjUEFRPT06fDo6fDp2bFlyazVWNmdCaU1BMTVOM0U4SGNBPT0=",
			want: 1,
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := New().Extract(tt.url)
			if err != nil {
				t.Error(err)
			}
			if len(data) > tt.want || len(data) == 0 {
				t.Errorf("Got: %v - want: %v", len(data), tt.want)
			}
		})
	}
}

func TestFindJWPlayerURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "HTML string",
			url:  `<iframe src="https://hentaihaven.xxx/wp-content/plugins/player-logic/player.php?data=OTZPUXF5OU1DbWxVSzROajdkWkx6T1Y4RFBVZ3JQKzRqWlB3M1FXY2tudUU2S1VHNklsUHY0TmRWREtCd3pKb2Q3ZVVtNklCWDJRUndqNGZQWjR5MjJKOEo2ZXU3QzBmTEU5aTRKK2tGME1MazZZUzh4Y0ZMRThyOG9TOTlHMnM0TnZvQnRIelhHTTJrL1htYUVjUEFRPT06fDo6fDp2bFlyazVWNmdCaU1BMTVOM0U4SGNBPT0=" frameborder="0" scrolling="no" allowfullscreen=""></iframe>`,
			want: "https://hentaihaven.xxx/wp-content/plugins/player-logic/player.php?data=OTZPUXF5OU1DbWxVSzROajdkWkx6T1Y4RFBVZ3JQKzRqWlB3M1FXY2tudUU2S1VHNklsUHY0TmRWREtCd3pKb2Q3ZVVtNklCWDJRUndqNGZQWjR5MjJKOEo2ZXU3QzBmTEU5aTRKK2tGME1MazZZUzh4Y0ZMRThyOG9TOTlHMnM0TnZvQnRIelhHTTJrL1htYUVjUEFRPT06fDo6fDp2bFlyazVWNmdCaU1BMTVOM0U4SGNBPT0=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := FindJWPlayerURL(&tt.url)
			if u == "" {
				t.Errorf("Got: %v - want: %v", u, tt.want)
			}
		})
	}
}
