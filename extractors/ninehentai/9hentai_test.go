package ninehentai

import "testing"

//9hentai
func TestParseURL(t *testing.T) {
	tests := []struct {
		Name string
		in   string
		Want int
	}{
		{
			Name: "Single Gallery",
			in:   "https://9hentai.to/g/301/",
			Want: 1,
		}, {
			Name: "Single Gallery .ru",
			in:   "https://www1.9hentai.ru/g/71163/",
			Want: 1,
		},
		{
			Name: "Single Tag",
			in:   "https://9hentai.to/t/71/",
			Want: 18,
		}, {
			Name: "Complex search",
			in:   "https://9hentai.to/t/71/#~(text~'~page~0~sort~0~pages~(range~(~0~2000))~tag~(text~'AN~type~1~tags~(~)~items~(included~(~(id~71~name~'Alice~description~null~type~5~books_count~25)~(id~30~name~'Anal~description~null~type~1))~excluded~(~))))#",
			Want: 18,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			galleries, err := parseURL(tt.in)
			if err != nil {
				t.Error(err)
			}
			if len(galleries) < tt.Want {
				t.Errorf("Got: %v - Want: %v", len(galleries), tt.Want)
			}
		})
	}
}

func TestExtract(t *testing.T) {
	tests := []struct {
		Name string
		in   string
		Want int
	}{
		{
			Name: "Single Gallery",
			in:   "https://9hentai.to/g/301/",
			Want: 1,
		}, {
			Name: "Complex search",
			in:   "https://9hentai.to/t/71/#~(text~'~page~0~sort~0~pages~(range~(~0~2000))~tag~(text~'AN~type~1~tags~(~)~items~(included~(~(id~71~name~'Alice~description~null~type~5~books_count~25)~(id~30~name~'Anal~description~null~type~1))~excluded~(~))))#",
			Want: 18,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			data, err := New().Extract(tt.in)
			if err != nil {
				t.Error(err)
			}
			if len(data) < tt.Want {
				t.Errorf("Got: %v - Want: %v", len(data), tt.Want)
			}
		})
	}
}
