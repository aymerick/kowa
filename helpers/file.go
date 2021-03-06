package helpers

import (
	"fmt"
	"log"
	"os"
	"path"
)

// FileBase returns base file name
// Example: /foo/bar/baz.png => baz
func FileBase(filePath string) string {
	fileName := path.Base(filePath)
	fileExt := path.Ext(filePath)

	return fileName[:len(fileName)-len(fileExt)]
}

// AvailableFilePath returns filePath is available, otherwise returns a path with a random filename
func AvailableFilePath(filePath string) string {
	// check if file does not exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return filePath
	}

	// generate a new file name
	fileBaseName := FileBase(filePath)
	fileExt := path.Ext(filePath)
	fileDir := path.Dir(filePath)

	newFilePath := path.Join(fileDir, fmt.Sprintf("%s_%s%s", fileBaseName, RandomAlphaNumString(16), fileExt))

	// check if file already exists
	if _, err := os.Stat(newFilePath); !os.IsNotExist(err) {
		panic(fmt.Sprintf("Random file already exists: %s", newFilePath))
	}

	return newFilePath
}

// EnsureDirectory ensures that given directory exists
func EnsureDirectory(dirPath string) {
	log.Printf("[DBG] Ensuring dir: %s", dirPath)

	err := os.MkdirAll(dirPath, 0777)
	if err != nil && err != os.ErrExist {
		panic(err)
	}
}

// WorkingDir returns current working directory
func WorkingDir() string {
	result, _ := os.Getwd()

	return result
}
