package static

type URL struct {
	URL string
	Ext string
}

type Stream struct {
	URLs    []URL  `json: "url"`
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
