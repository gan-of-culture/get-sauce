package static

// URL Struct of URL
type URL struct {
	URL string
	Ext string
}

// Stream Struct of stream
type Stream struct {
	URLs    []URL  `json:"url"`
	Quality string `json:"quality"`
	Size    int64  `json:"size"`
	Info    string `json:"info"`
}

// Data Struct of data
type Data struct {
	Site  string `json:"site"`
	Title string `json:"title"`
	Type  string `json:"type"`

	Streams map[string]Stream `json:"streams"`

	Url string `json:"sourceUrl"`
}
