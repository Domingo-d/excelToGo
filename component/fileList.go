package component

import (
	"excelToGo/common"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"path/filepath"
)

func BuildFileList(data *common.AppData) fyne.CanvasObject {
	data.ListWidget = widget.NewList(
		func() int { return len(data.FilteredList) },
		func() fyne.CanvasObject { return widget.NewLabel("模板") },
		func(i int, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(filepath.Base(data.FilteredList[i]))
		},
	)

	return data.ListWidget
}
