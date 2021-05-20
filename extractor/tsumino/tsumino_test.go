package tsumino

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Gallery",
			url:  "https://www.tsumino.com/entry/55285",
			want: 1,
		}, {
			name: "Tag",
			url:  "https://www.tsumino.com/books#~(Sort~'Rating~Include~(~)~Tags~(~(Type~1~Text~'Angel~Exclude~false)))#",
			want: 36,
		}, {
			name: "Tag different page",
			url:  "https://www.tsumino.com/books#~(PageNumber~3~Sort~'Rating~Include~(~)~Tags~(~(Type~1~Text~'Angel~Exclude~false)))#",
			want: 36,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urls := ParseURL(tt.url)
			if len(urls) > tt.want || len(urls) == 0 {
				t.Errorf("Got: %v - want: %v", len(urls), tt.want)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Gallery",
			url:  "https://www.tsumino.com/entry/55285",
			want: 1,
		}, /*{
			name: "Tag different page",
			url:  "https://www.tsumino.com/books#~(PageNumber~3~Sort~'Rating~Include~(~)~Tags~(~(Type~1~Text~'Angel~Exclude~false)))#",
			want: 25,
		},*/ //somtimes you run into a captcha
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Extract(tt.url)
			if err != nil {
				t.Error(err)
			}
			if len(data) > tt.want || len(data) == 0 {
				t.Errorf("Got: %v - want: %v", len(data), tt.want)
			}
		})
	}
}
