package merger

import (
	"fmt"
	"os"

	"github.com/gan-of-culture/get-sauce/utils"
	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
)

// fragmentMerger merge stream fragments into one final file e.g M3U parts.
type fragmentMerger struct {
	decryptKey  []byte
	bar         bool
	progressBar *progressbar.ProgressBar
}

func (fM *fragmentMerger) Merge(files []*MergeFile, outFile string) error {
	lenOfFiles := len(files)
	if lenOfFiles <= 1 {
		return nil
	}

	file, err := os.OpenFile(outFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	fM.progressBar = utils.InitPB(utils.ProgressBarConfig{
		Length:      int64(lenOfFiles),
		Description: fmt.Sprintf("Merging into %s ...", file.Name()),
		AsBytes:     false,
	})

	var d []byte
	for _, f := range files {
		if len(fM.decryptKey) > 0 {
			d, err = decrypt(fM.decryptKey, f.Path)
			if err != nil {
				return err
			}
		} else {
			d, err = os.ReadFile(f.Path)
			if err != nil {
				return errors.WithStack(err)
			}
		}

		if _, err := file.Write(d); err != nil {
			return errors.WithStack(err)
		}

		if fM.bar {
			fM.progressBar.Add(1)
		}
	}

	return nil
}

func NewFragmentMerger(decryptKey []byte, bar bool) merger {
	return &fragmentMerger{decryptKey: decryptKey, bar: bar}
}
