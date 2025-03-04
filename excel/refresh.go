package excel

import (
	"excelToGo/common"
	"github.com/fsnotify/fsnotify"
	"log"
	"path/filepath"
)

func RefreshFileList(data *common.AppData) {
	data.FileList, _ = filepath.Glob(filepath.Join(data.CurrentDir, "*.xlsx"))
	data.FilteredList = data.FileList
	data.ListWidget.Refresh()
}

func StartWatcher(data *common.AppData) {
	watcher, _ := fsnotify.NewWatcher()
	watcher.Add(data.CurrentDir)

	go func() {
		for {
			select {
			case <-watcher.Events:
				RefreshFileList(data)
			case err := <-watcher.Errors:
				if nil != err {
					log.Println("监控错误:", err)
					panic(err)
				}
			}
		}
	}()
}
