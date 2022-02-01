package config

var (
	//Amount of files to download
	Amount int
	// Caption to download if available
	Caption int
	// Keep video, audio and subtitles. Don't merge using ffmpeg
	Keep bool
	// OutputPath for files
	OutputPath string
	// OutputName for file
	OutputName string
	// Pages of doujinshi
	Pages string
	// Quiet mode - show minimal information
	Quiet bool
	// SelectStream to download
	SelectStream string
	// ShowExtractedData as json
	ShowExtractedData bool
	// ShowInfo of all available streams
	ShowInfo bool
	// Truncate file if it already exists
	Truncate bool
	// Workers for downloading
	Workers int
	// Username for exhentai.org
	Username string
	// UserPassword for exhentai.org
	UserPassword string
)

// FakeHeaders fake http headers
var FakeHeaders = map[string]string{
	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	"Accept-Language": "en-US,en;q=0.8",
	"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.81 Safari/537.36",
}
