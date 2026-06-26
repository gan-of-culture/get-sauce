package merger

import (
	"archive/zip"
	"fmt"
	"os"
	"path"

	"github.com/gan-of-culture/get-sauce/static"
	"github.com/gan-of-culture/get-sauce/utils"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
)

// archiveMerger merge stream files into one file archive file e.g a cbz archive.
type archiveMerger struct {
	data        *static.Data
	bar         bool
	progressBar *progressbar.ProgressBar
}

func (aM *archiveMerger) Merge(files []*MergeFile, outFile string) error {

	file, err := os.OpenFile(outFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	w := zip.NewWriter(file)

	comicInfo, err := NewComicInfo(aM.data)
	if err != nil {
		return err
	}

	f, err := w.Create("ComicInfo.xml")
	if err != nil {
		return errors.WithStack(err)
	}
	if _, err = f.Write(comicInfo); err != nil {
		return errors.WithStack(err)
	}

	aM.progressBar = utils.InitPB(utils.ProgressBarConfig{
		Length:      int64(len(files)),
		Description: fmt.Sprintf("Merging into %s ...", file.Name()),
		AsBytes:     false,
	})

	for _, file := range files {
		_, fname := path.Split(file.Path)

		body, err := os.ReadFile(file.Path)
		if err != nil {
			return errors.WithStack(err)
		}

		f, err := w.Create(fname)
		if err != nil {
			return errors.WithStack(err)
		}
		if _, err = f.Write([]byte(body)); err != nil {
			return errors.WithStack(err)
		}

		if aM.bar {
			aM.progressBar.Add(1)
		}
	}

	if err := w.Close(); err != nil {
		return err
	}

	// delete only after the zip has been created without errors
	for _, f := range files {
		err = os.Remove(f.Path)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewArchiveMerger(bar bool, data *static.Data) merger {
	return &archiveMerger{data: data, bar: bar}
}
