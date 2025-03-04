package component

import (
	"excelToGo/common"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"strings"
)

func BuildSearchBar(data *common.AppData) fyne.CanvasObject {
	data.SearchEntry.OnChanged = func(s string) {
		data.FilteredList = filterFiles(data.FileList, s)
		data.ListWidget.Refresh()
	}

	return container.NewBorder(nil, nil, widget.NewLabel("搜索:"), data.SearchEntry)
}

func filterFiles(files []string, keyword string) []string {
	var result []string
	lowerKeyword := strings.ToLower(keyword)
	for _, f := range files {
		if strings.Contains(strings.ToLower(f), lowerKeyword) {
			result = append(result, f)
		}
	}

	return result
}
