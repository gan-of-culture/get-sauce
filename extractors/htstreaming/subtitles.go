package htstreaming

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

type caption struct {
	Kind     string `json:"kind"`
	File     string `json:"file"`
	Label    string `json:"label"`
	Language string `json:"language"`
	Default  bool   `json:"default"`
}

var reSubtileParse = regexp.MustCompile(`\b\w+\b`)
var reSubtitleParams = regexp.MustCompile(`'([^']+)',(\d+),(\d+),'([^']+)`)
var reSubtitles = regexp.MustCompile(`{"kind":"captions"[^}]*}`)

func parseCaptions(jsParams string) []*static.Caption {
	captions := []*static.Caption{}

	for _, c := range reSubtitles.FindAllString(jsParams, -1) {
		caption := caption{}
		err := json.Unmarshal([]byte(c), &caption)
		if err != nil {
			continue
		}
		if caption.Language == "" {
			caption.Language = caption.Label
		}

		captions = append(captions, &static.Caption{
			URL: static.URL{
				URL: caption.File,
				Ext: utils.GetFileExt(caption.File),
			},
			Language: caption.Language,
		})
	}

	return captions
}

func parseFirePlayerParams(jsTemplate string, a, c int, keywords []string) string {
	dict := map[string]string{}

	c -= 1
	for ; c > -1; c-- {
		if keywords[c] == "" {
			keywords[c] = e(c, &a)
		}
		dict[e(c, &a)] = keywords[c]
	}

	jsDict := reSubtileParse.ReplaceAllStringFunc(jsTemplate, func(s string) string { return dict[s] })
	jsDict = strings.ReplaceAll(jsDict, `\\`, `\`)
	return jsDict
}

func e(c int, a *int) string {
	val1 := ""
	if c >= *a {
		val1 = e(c / *a, a)
	}
	val2 := ""
	c = c % *a
	if c <= 35 {
		val2 = strconv.FormatInt(int64(c), 36)
	} else {
		val2 = string(rune(c + 29))
	}

	return val1 + val2
}
