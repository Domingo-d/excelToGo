package common

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type (
	AppData struct {
		CurrentDir   string
		FileList     []string
		FilteredList []string
		SearchEntry  *widget.Entry
		ListWidget   *widget.List
		RightPanel   *fyne.Container
	}
)
