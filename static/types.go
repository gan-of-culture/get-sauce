package static

// URL Struct of URL
type URL struct {
	// URL that contains the data to be downloaded
	URL string `json:"url"`
	// Ext of the data
	Ext string `json:"ext"`
}

// Stream Struct of stream
type Stream struct {
	URLs []*URL `json:"url"`
	// Quality e.g. 2160p, 1080p, 720p ... or 1050x1200 or codec
	Quality string `json:"quality"`
	// Size of stream - fill this if it is a big blob of data and you want to make use of concurrent downloading
	Size int64 `json:"size"`
	// Info that could be interesting for the user
	Info string `json:"info"`
	// Ext after the files are merged
	Ext string `json:"ext"`
	// Key that is needed to decrypt this stream
	Key []byte `json:"key"`
}

// DataType indicates the type of extracted data, e.g. video or image.
type DataType string

const (
	// DataTypeVideo indicates the type of extracted data is the video.
	DataTypeVideo DataType = "video"
	// DataTypeImage indicates the type of extracted data is the image.
	DataTypeImage DataType = "image"
	// DataTypeImage indicates the type of extracted data is the unknown.
	DataTypeUnknown DataType = "unknown"
)

// Data Struct of data
type Data struct {
	// Site name of the media host
	Site string `json:"site"`
	// Title of data
	Title string `json:"title"`
	// Type of data commonly image or video
	Type DataType `json:"type"`

	// Streams of different quality or mirrors
	Streams map[string]*Stream `json:"streams"`

	// Url that was supplied to the scraper
	Url string `json:"sourceUrl"`
}

type Extractor interface {
	Extract(URL string) ([]*Data, error)
}
