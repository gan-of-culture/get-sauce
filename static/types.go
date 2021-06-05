package static

// URL Struct of URL
type URL struct {
	// URL that contains the data to be downloaded
	URL string
	// Ext of the data
	Ext string
}

// Stream Struct of stream
type Stream struct {
	URLs []URL `json:"url"`
	// Quality e.g. 2160p, 1080p, 720p ... or 1050x1200 or codec
	Quality string `json:"quality"`
	// Size of stream - fill this if it is a big blob of data and you want to make use of concurrent downloading
	Size int64 `json:"size"`
	// Info that could be interesting for the user
	Info string `json:"info"`
}

// Data Struct of data
type Data struct {
	// Site name of the media host
	Site string `json:"site"`
	// Title of data
	Title string `json:"title"`
	// Type of data commonly image or video -> but if needed or possible it's the mimeType of the data
	Type string `json:"type"`

	// Streams of different quality or mirrors
	Streams map[string]Stream `json:"streams"`

	// Url that was supplied to the scraper
	Url string `json:"sourceUrl"`
}

type Extractor interface {
	Extract(URL string) ([]*Data, error)
}
