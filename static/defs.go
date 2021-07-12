package static

import (
	"errors"
)

var (
	// ErrURLParseFailed defines URL parse failed error.
	ErrURLParseFailed = errors.New("URL parse failed")
	// ErrLoginRequired defines login required error.
	ErrLoginRequired = errors.New("login required")
	// ErrDataSourceParseFailed defines a data source parse error.
	ErrDataSourceParseFailed = errors.New("data source parse failed")
)
