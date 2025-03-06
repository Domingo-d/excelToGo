package cmd

import (
	"bufio"
	"excelToGo/service"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"time"
)

var (
	InteractiveCmd = &cobra.Command{
		Use:               "Console",
		Short:             "交互式控制台",
		PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}

	quit            chan os.Signal
	wait            *sync.WaitGroup
	interactiveQuit chan bool
	scanner         *bufio.Reader
	red             = color.New(color.FgRed)
	green           = color.New(color.FgGreen)
)

func run() {
	quit = make(chan os.Signal, 1)
	wait = &sync.WaitGroup{}
	interactiveQuit = make(chan bool, 1)
	scanner = bufio.NewReader(os.Stdin)

	wait.Add(1)
	go func() {
		service.InitFileList()
		wait.Done()
	}()

	wait.Add(1)
	go interactive()

	signal.Notify(quit, os.Interrupt)
	<-quit
	close(quit)

	if interactiveQuit != nil {
		interactiveQuit <- true
	}

	service.CloseFileList()
	wait.Wait()

	color.Green("%s, Shutdown Server ...\r\n", time.Now().Format("2006-01-02 15:04:05"))
}

func interactive() {
outLoop:
	for {
		select {
		case <-interactiveQuit:
			break outLoop
		default:
			showMenu()
			green.Print("> exit退出 输入: ")
			inText, err := scanner.ReadString('\n')
			if err != nil {
				color.Red("scanner.ReadLine Error:%s", err.Error())
				continue
			}

			cmd := strings.TrimSpace(inText)
			if cmd == "exit" {
				break outLoop
			}

			handleInput(cmd)
		}
	}

	close(interactiveQuit)
	interactiveQuit = nil
	wait.Done()
	quit <- os.Interrupt
}

func showMenu() {
	menu := `
	============== Excel to Go ==============
	1. 列出当前目录下的所有的excel文件
	2. 查询文件
	3. 导出配置
	4. 转换go struct
	0. 重新显示菜单
	============== End ==============`

	color.Yellow(menu)
}

func handleInput(input string) {
	switch input {
	case "1", "list":
		service.ListFiles()
	case "2", "search":
		searchFile()
	case "3", "generate":
		generateJson()
	case "4", "struct":
		excelToStruct()
	}
}

func searchFile() {
	green.Print("> ")
	inText, err := scanner.ReadString('\n')
	if err != nil {
		color.Red("searchFile scanner.ReadString Error:%s", err.Error())
		return
	}
	fileName := strings.TrimSpace(inText)
	service.SearchFile(fileName)
}

func generateJson() {
	closeWriteList := make([]io.WriteCloser, 0)
	closeReadList := make([]io.ReadCloser, 0)

	defer func() {
		for _, v := range closeWriteList {
			v.Close()
		}

		for _, v := range closeReadList {
			v.Close()
		}
	}()

	for {
		green.Print("> exit退出 输入文件名:")
		inText, err := scanner.ReadString('\n')
		if err != nil {
			color.Red("generateJson scanner.ReadString Error:%s", err.Error())
			return
		}

		if strings.Contains(inText, "exit") {
			return
		}

		filePath := service.GetPath(strings.TrimSpace(inText))
		if filePath == "" {
			color.Red("generateJson filePath is empty")
			return
		}

		cmd := exec.Command("./导出配置.exe", filePath)
		stdin, err := cmd.StdinPipe()
		stdout, err1 := cmd.StdoutPipe()
		if err != nil || err1 != nil {
			color.Red("获取stdinPipe失败:%v", err)
		}

		closeWriteList = append(closeWriteList, stdin)
		closeReadList = append(closeReadList, stdout)

		if err := cmd.Start(); err != nil {
			color.Red("generateJson cmd.Start Error:%s", err.Error())
			return
		}

		//go func() {
		if _, err := io.WriteString(stdin, filePath+"\n"); err != nil {
			color.Red("stdin.Write Error:%s", err.Error())
			return
		}
		//}()

		go func() {
			//outText, err := io.ReadAll(stdout)
			//if err != nil {
			//	color.Red("io.ReadAll Error:%s", err.Error())
			//}

			var outText string
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				outText = scanner.Text()
				color.Green("%s", outText)

				if strings.Contains(outText, "导表完成") {
					io.WriteString(stdin, "\n")
				}
			}

			color.Green("%s", outText)
		}()

		color.Yellow("pid:%d", cmd.Process.Pid)
		if err := cmd.Wait(); err != nil {
			color.Red("generateJson cmd.Wait Error:%s", err.Error())
		}

		//output, err := cmd.CombinedOutput()
		//if err != nil {
		//	color.Red("generateJson cmd.CombinedOutput Error:%s", err.Error())
		//}
	}
}

func excelToStruct() {
	for {
		green.Print("> exit退出 输入:")
		inText, err := scanner.ReadString('\n')
		if err != nil {
			color.Red("generateJson scanner.ReadString Error:%s", err.Error())
			return
		}

		if strings.Contains(inText, "exit") {
			return
		}

		filePath := service.GetPath(strings.TrimSpace(inText))
		if filePath == "" {
			color.Red("generateJson filePath is empty")
			return
		}

		service.ExcelToGo(filePath)
	}
}
