package ohentai

import "testing"

func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want int
	}{
		{
			name: "Single Video",
			url:  "https://ohentai.org/detail.php?vid=MzM5MQ==",
			want: 1,
		}, {
			name: "Tag",
			url:  "https://ohentai.org/tagsearch.php?tag=Uncensored",
			want: 24,
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
			name: "Single Video",
			url:  "https://ohentai.org/detail.php?vid=MzM5MQ==",
			want: 1,
		}, {
			name: "Tag",
			url:  "https://ohentai.org/tagsearch.php?tag=Uncensored",
			want: 24,
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
