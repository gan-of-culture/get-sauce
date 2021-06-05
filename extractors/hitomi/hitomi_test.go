package hitomi

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Gallery",
			url:  "https://hitomi.la/doujinshi/%E3%82%B8%E3%82%A7%E3%83%B3%E3%83%88%E3%83%AB%E3%82%B3%E3%83%8D%E3%82%AF%E3%83%88!re:dive-%E6%97%A5%E6%9C%AC%E8%AA%9E-1905632.html",
			want: 1,
		}, {
			name: "Tag",
			url:  "https://hitomi.la/tag/artbook-all.html",
			want: 25,
		}, {
			name: "Tag different page",
			url:  "https://hitomi.la/tag/artbook-all.html?page=2",
			want: 25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urls := parseURL(tt.url)
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
			url:  "https://hitomi.la/doujinshi/%E3%82%B8%E3%82%A7%E3%83%B3%E3%83%88%E3%83%AB%E3%82%B3%E3%83%8D%E3%82%AF%E3%83%88!re:dive-%E6%97%A5%E6%9C%AC%E8%AA%9E-1905632.html",
			want: 1,
		}, {
			name: "Single Gallery",
			url:  "https://hitomi.la/manga/%E7%8C%A5%E8%A4%BB%E3%83%9F%E3%82%B5%E3%82%A4%E3%83%AB-%E6%97%A5%E6%9C%AC%E8%AA%9E-440479.html",
			want: 1,
		},
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
