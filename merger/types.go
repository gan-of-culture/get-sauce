package merger

import (
	"github.com/gan-of-culture/get-sauce/static"
)

type MergeFile struct {
	Path     string
	DataType static.DataType
}

type merger interface {
	Merge(files []*MergeFile, outFile string) error
}
