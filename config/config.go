package config

var (
	// Pages of doushin
	Pages string
	// OutputPath
	OutputPath string
	// OutputName
	OutputName string
	// ShowInfo
	ShowInfo bool
	// SelectStream
	SelectStream string
)

// FakeHeaders fake http headers
var FakeHeaders = map[string]string{
	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	"Accept-Language": "en-US,en;q=0.8",
	"Accept-Encoding": "gzip, deflate, br",
	"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.81 Safari/537.36",
}
