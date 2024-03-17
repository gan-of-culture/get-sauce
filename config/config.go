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
	// Timeout for the http.client in minutes
	Timeout int
	// Truncate file if it already exists
	Truncate bool
	// UserHeaders for the HTTP requests. To bypass Cloudflare or DDOS-GUARD protection
	UserHeaders string
	// Workers for downloading
	Workers int
)

// FakeHeaders fake http headers
var FakeHeaders = map[string]string{
	"Accept":          "*/*",
	"Accept-Language": "en-US,en;q=0.8",
	"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.81 Safari/537.36",
}

var FakeHeadersFirefox117 = map[string]string{
	"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/117.0",
	"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8",
	"Accept-Language":           "en-US,en;q=0.5",
	"Accept-Encoding":           "gzip, deflate, br",
	"Upgrade-Insecure-Requests": "1",
	"Sec-Fetch-Dest":            "document",
	"Sec-Fetch-Mode":            "navigate",
	"Sec-Fetch-Site":            "none",
	"Sec-Fetch-User":            "?1",
	"TE":                        "Trailers",
	"X-Requested-With":          "XMLHttpRequest",
}
