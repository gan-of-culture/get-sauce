package hentaiplay

import (
	"encoding/json"
	"net/url"
	"regexp"
	"slices"
	"strings"

	"github.com/gan-of-culture/get-sauce/request"
	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
)

type extractor struct{}

type post struct {
	ID      int    `json:"id"`
	Date    string `json:"date"`
	DateGmt string `json:"date_gmt"`
	GUID    struct {
		Rendered string `json:"rendered"`
	} `json:"guid"`
	Modified    string `json:"modified"`
	ModifiedGmt string `json:"modified_gmt"`
	Slug        string `json:"slug"`
	Status      string `json:"status"`
	Type        string `json:"type"`
	Link        string `json:"link"`
	Title       struct {
		Rendered string `json:"rendered"`
	} `json:"title"`
	Content struct {
		Rendered  string `json:"rendered"`
		Protected bool   `json:"protected"`
	} `json:"content"`
	Excerpt struct {
		Rendered  string `json:"rendered"`
		Protected bool   `json:"protected"`
	} `json:"excerpt"`
	Author        int    `json:"author"`
	FeaturedMedia int    `json:"featured_media"`
	CommentStatus string `json:"comment_status"`
	PingStatus    string `json:"ping_status"`
	Sticky        bool   `json:"sticky"`
	Template      string `json:"template"`
	Format        string `json:"format"`
	Meta          struct {
		Footnotes string `json:"footnotes"`
	} `json:"meta"`
	Categories []int `json:"categories"`
	Tags       []int `json:"tags"`
}

const site = "https://hentaiplay.net/"
const API = "https://hentaiplay.net/wp-json/wp/v2"

var reVideoURL = regexp.MustCompile(`clip-link[\s\S]*?href="([^"]+)`)
var reSource = regexp.MustCompile(`<source src="([^"]+)`)

func parseURL(URL string) ([]string, error) {
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	slug := slices.DeleteFunc(strings.Split(u.Path, "/"), func(seg string) bool { return seg == "" })[0]

	APIURL, err := url.Parse(API)
	if err != nil {
		return nil, err
	}
	APIURL.Path, err = url.JoinPath(APIURL.Path, "posts")
	if err != nil {
		return nil, err
	}
	q := APIURL.Query()
	q.Add("slug", slug)
	APIURL.RawQuery = q.Encode()

	JSONBytes, err := request.GetAsBytes(APIURL.String())
	if err != nil {
		return nil, err
	}

	var posts []post
	err = json.Unmarshal(JSONBytes, &posts)
	if err != nil {
		return nil, err
	}
	if len(posts) > 0 {
		return []string{URL}, nil
	}

	body, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	body = strings.Split(body, `id="sidebar"`)[0]

	var out []string
	for _, URLPart := range reVideoURL.FindAllStringSubmatch(body, -1) {
		out = append(out, utils.GetLastItemString(URLPart))
	}

	return out, nil
}

func extractData(URL string) (*static.Data, error) {
	body, err := request.Get(URL)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	slug := slices.DeleteFunc(strings.Split(u.Path, "/"), func(seg string) bool { return seg == "" })[0]

	APIURL, err := url.Parse(API)
	if err != nil {
		return nil, err
	}
	APIURL.Path, err = url.JoinPath(APIURL.Path, "posts")
	if err != nil {
		return nil, err
	}
	q := APIURL.Query()
	q.Add("slug", slug)
	APIURL.RawQuery = q.Encode()

	JSONBytes, err := request.GetAsBytes(APIURL.String())
	if err != nil {
		return nil, err
	}

	var posts []post
	err = json.Unmarshal(JSONBytes, &posts)
	if err != nil {
		return nil, err
	}
	if len(posts) < 1 {
		return nil, static.ErrDataSourceParseFailed
	}
	post := posts[0]
	sourceURL := utils.GetLastItemString(reSource.FindStringSubmatch(body))

	size, err := request.Size(sourceURL, URL)
	if err != nil {
		return nil, err
	}

	return &static.Data{
		Site:  site,
		Title: post.Title.Rendered,
		Type:  static.DataTypeVideo,
		Streams: map[string]*static.Stream{
			"0": {
				Type: static.DataTypeVideo,
				URLs: []*static.URL{{URL: sourceURL, Ext: utils.GetFileExt(sourceURL)}},
				Size: size,
			},
		},
		URL: URL,
	}, nil
}

// Extract implements [static.Extractor].
func (e extractor) Extract(URL string) ([]*static.Data, error) {
	URLs, err := parseURL(URL)
	if err != nil {
		return nil, err
	}
	if len(URLs) == 0 {
		return nil, static.ErrURLParseFailed
	}

	data := []*static.Data{}
	for _, u := range URLs {
		d, err := extractData(u)
		if err != nil {
			return nil, utils.Wrap(err, u)
		}
		data = append(data, d)
	}

	return data, nil
}

// New returns a hentaiplay extractor
func New() static.Extractor {
	return extractor{}
}
