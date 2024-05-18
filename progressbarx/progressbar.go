package progressbarx

import "github.com/schollz/progressbar/v3"

type ProgressBar struct {
	*progressbar.ProgressBar
}

func NewProgressBar(total int, options ...progressbar.Option) *ProgressBar {
	return &ProgressBar{progressbar.NewOptions(total)}
}
