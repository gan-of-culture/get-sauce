package static

type Stream struct {
	Url     string `json: "url"`
	Quality string `json: "quality"`
	Size    int64  `json: "size"`
}

type Data struct {
	Site  string `json: "site"`
	Title string `json: "title"`
	Type  string `json: "type"`

	Streams map[string]Stream `json: "streams"`

	Err error  `json: "err"`
	Url string `json: "sourceUrl"`
}
