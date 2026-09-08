package main

import (
	"os"
	"path/filepath"
	"time"
)

const (
	updateCacheDirectory = "UpdateCache"
	updateFileName       = "PoeBuy.exe"
)

func main() {

	// Get the path of the currently running executable
	exePath, err := os.Executable()
	if err != nil {
		return
	}

	currentDir := filepath.Dir(exePath)

	// Check if the current directory is the UpdateCache directory
	if filepath.Base(currentDir) != updateCacheDirectory {
		return
	}

	sourceFile := filepath.Join(currentDir, updateFileName)
	targetFile := filepath.Join(filepath.Dir(currentDir), updateFileName)

	// Retry the rename operation multiple times in case the file is still in use
	for i := 0; i < 10; i++ {
		err := os.Rename(sourceFile, targetFile)
		if err == nil {
			return
		}

		time.Sleep(500 * time.Millisecond)
	}
}
