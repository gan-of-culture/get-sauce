package test

import (
	"encoding/json"
	"testing"

	"github.com/gan-of-culture/get-sauce/static"
)

// Args Arguments for extractor tests
type Args struct {
	URL     string
	Title   string
	Quality string
	Size    int64
}

// CheckData check the given data
func CheckData(args, data Args) bool {
	if args.Title != data.Title {
		return false
	}
	// not every video got quality information
	if data.Quality != "" && args.Quality != data.Quality {
		return false
	}
	if data.Size != 0 && args.Size != data.Size {
		return false
	}
	return true
}

// Check check the result
func Check(t *testing.T, args Args, data *static.Data) {
	defaultData := data.Streams["0"]

	if defaultData == nil {
		t.Errorf("Data contains no streams or no default stream")
		return
	}

	temp := Args{
		Title:   data.Title,
		Quality: defaultData.Quality,
		Size:    defaultData.Size,
	}
	if !CheckData(args, temp) {
		jsonData, _ := json.MarshalIndent(defaultData, "", "    ")
		t.Log(jsonData)
		t.Errorf("Got: %v\nExpected: %v", temp, args)
	}
}

// CheckError check the error
func CheckError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
