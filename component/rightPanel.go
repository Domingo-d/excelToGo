package component

import (
	"excelToGo/common"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func BuildRightPanel(data *common.AppData) fyne.CanvasObject {
	return container.NewVBox(
		widget.NewButton("打开文件", func() {}),
		widget.NewButton("导出", func() {

		}),

		widget.NewButton("生成Struct", func() {}),
	)
}
