package main

import (
	"excelToGo/common"
	"excelToGo/component"
	"excelToGo/excel"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	window := myApp.NewWindow("Excel to Go v1.0")
	data := &common.AppData{CurrentDir: ".", SearchEntry: widget.NewEntry()}
	window.SetContent(container.NewBorder(
		component.BuildSearchBar(data),
		nil,
		component.BuildFileList(data),
		component.BuildRightPanel(data),
	))

	excel.RefreshFileList(data)
	excel.StartWatcher(data)
	window.ShowAndRun()

	//fileName := flag.String("file", "data.xlsx", "file to read")
	//flag.Parse()
	//
	//err := excel.ExcelToGo(*fileName)
	//if nil != err {
	//	log.Fatalf("Error generating Go file: %s", err)
	//}
}
