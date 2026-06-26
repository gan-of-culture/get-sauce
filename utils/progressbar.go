package utils

import "github.com/schollz/progressbar/v3"

type ProgressBarConfig struct {
	Length      int64
	Description string
	AsBytes     bool
}

func InitPB(config ProgressBarConfig) *progressbar.ProgressBar {

	if config.AsBytes {
		return progressbar.NewOptions(
			int(config.Length),
			progressbar.OptionShowBytes(true),
			progressbar.OptionSetDescription(config.Description),
			progressbar.OptionSetPredictTime(true),
			progressbar.OptionSetRenderBlankState(true),
		)
	}

	return progressbar.NewOptions(
		int(config.Length),
		progressbar.OptionShowIts(),
		progressbar.OptionSetDescription(config.Description),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetRenderBlankState(true),
	)
}
