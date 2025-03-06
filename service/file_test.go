package service

import (
	"sync"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	wait := sync.WaitGroup{}

	wait.Add(1)
	go func() {
		InitFileList()
		wait.Done()
	}()

	time.Sleep(10 * time.Second)
	ListFiles()

	CloseFileList()

	wait.Wait()
}
