package ninehentai

import "testing"

//9hentai
func TestParseURL(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{
			name: "Single Gallery",
			in:   "https://9hentai.to/g/301/",
			want: 1,
		}, {
			name: "Single Gallery .ru",
			in:   "https://www1.9hentai.ru/g/71163/",
			want: 1,
		},
		{
			name: "Single Tag",
			in:   "https://9hentai.to/t/71/",
			want: 18,
		}, {
			name: "Complex search",
			in:   "https://9hentai.to/t/71/#~(text~'~page~0~sort~0~pages~(range~(~0~2000))~tag~(text~'AN~type~1~tags~(~)~items~(included~(~(id~71~name~'Alice~description~null~type~5~books_count~25)~(id~30~name~'Anal~description~null~type~1))~excluded~(~))))#",
			want: 18,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			galleries, err := parseURL(tt.in)
			if err != nil {
				t.Error(err)
			}
			if len(galleries) < tt.want {
				t.Errorf("Got: %v - want: %v", len(galleries), tt.want)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{
			name: "Single Gallery",
			in:   "https://9hentai.to/g/301/",
			want: 1,
		}, {
			name: "Complex search",
			in:   "https://9hentai.to/t/71/#~(text~'~page~0~sort~0~pages~(range~(~0~2000))~tag~(text~'AN~type~1~tags~(~)~items~(included~(~(id~71~name~'Alice~description~null~type~5~books_count~25)~(id~30~name~'Anal~description~null~type~1))~excluded~(~))))#",
			want: 18,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := New().Extract(tt.in)
			if err != nil {
				t.Error(err)
			}
			if len(data) < tt.want {
				t.Errorf("Got: %v - want: %v", len(data), tt.want)
			}
		})
	}
}
