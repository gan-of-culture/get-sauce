package config

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

type MergeOpt string

const (
	MergeOptDefault MergeOpt = "default"
	MergeOptNone    MergeOpt = "none"
	MergeOptCBZ     MergeOpt = "cbz"
)

var mergeOpts = map[MergeOpt]string{
	MergeOptDefault: "default",
	MergeOptNone:    "none",
	MergeOptCBZ:     "cbz",
}

// String returns the string representation (required by flag.Value)
func (m MergeOpt) String() string {
	return mergeOpts[m]
}

// Set validates and sets the value (required by flag.Value)
func (m *MergeOpt) Set(value string) error {
	_, ok := mergeOpts[MergeOpt(value)]
	if !ok {
		return fmt.Errorf("invalid MergeOpt %q. Allowed values: %s", value, strings.Join(slices.Collect(maps.Values(mergeOpts)), ","))
	}
	*m = MergeOpt(value)
	return nil
}

var (
	// Amount of files to download
	Amount int
	// Caption to download if available
	Caption int
	// File download URLs listed in this file path
	File string
	// Merge output (default, none, cbz). CBZ only works if output is a stream of datatype image
	Merge MergeOpt
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
	// Subdirectory for the downloaded content. The directory name defaults to a cleaned up version of the data title
	Subdirectory bool
	// Timeout for the http.client in minutes
	Timeout int
	// Truncate file if it already exists
	Truncate bool
	// UserHeaders for the HTTP requests. To bypass Cloudflare or DDOS-GUARD protection
	UserHeaders string
	// Version prints current application version
	Version bool
	// Workers for downloading
	Workers int
)

// FakeHeaders fake http headers
var FakeHeaders = map[string]string{
	"Accept":          "*/*",
	"Accept-Language": "en-US,en;q=0.8",
	"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.3",
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
