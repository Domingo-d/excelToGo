package service

import (
	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	dir      = "E:/cqsy/envir/trunk/excels/"
	FileList map[string]string

	quit    chan bool
	watcher *fsnotify.Watcher
)

func GetPath(fileName string) string {
	return FileList[fileName]
}

func CloseFileList() {
	quit <- true
}

func InitFileList() {
	FileList = make(map[string]string)
	quit = make(chan bool, 1)

	loadFiles()

	tw, err := fsnotify.NewWatcher()
	if nil == err {
		watcher = tw

		dir, err := filepath.Abs(dir)
		if nil != err {
			color.Red("添加监控目录失败:", err, dir)
		} else {
			err = watcher.Add(dir)
		}

		runWatcher()

		watcher.Close()
	} else {
		color.Red("创建监控器失败:", err)
	}
}

func runWatcher() {
	count := 0
outLoop:
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				color.Red("文件监听服务退出")
				break
			}

			fileName := filepath.Base(event.Name)

			if strings.HasSuffix(fileName, ".xlsx") || strings.HasSuffix(fileName, ".xls") {
				switch event.Op {
				case fsnotify.Remove:
					delete(FileList, fileName)
				case fsnotify.Rename, fsnotify.Create:
					FileList[fileName] = event.Name
				}
			}

		case _ = <-quit:
			break outLoop

		default:
			count++
			//color.Green("runWatcher:%d", count)
			time.Sleep(time.Duration(5) * time.Second)
		}
	}

	color.Green("runWatcher:%d", count)
	close(quit)
}

func loadFiles() {
	err := filepath.WalkDir(dir, func(path string, info os.DirEntry, err error) error {
		if nil != err {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".xlsx" || ext == ".xls" {
			FileList[info.Name()] = path
		}

		return nil
	})

	if nil != err {
		color.Red("loadFiles Error:", err)
	}
}

func ListFiles() {
	for name, path := range FileList {
		color.Yellow("文件名: %s\t\t\t\t\t路径: %s\n", name, path)
	}
}

func SearchFile(fileName string) {
	for name, path := range FileList {
		if strings.Contains(name, fileName) {
			color.Yellow("文件名: %s\t\t\t\t\t路径: %s\n", name, path)
		}
	}
}
