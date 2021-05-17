package utils

import (
	"io/fs"
	"os"
	"path/filepath"
)

func GetFilenameSuffix(path string) string {
	return filepath.Ext(path)
}

func IsAbsolutePath(path string) bool {
	return filepath.IsAbs(path)
}

func JoinPath(args ...string) string {
	return filepath.Join(args...)
}

func SplitFilePath(path string) (string, string) {
	return filepath.Split(path)
}

func ExistsPath(path string) bool {
	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func IsDirectory(path string) bool {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func IsFile(path string) bool {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false
	}

	return !fileInfo.IsDir()
}

func CreatePaths(path string) bool {
	if ExistsPath(path) {
		return true
	}

	err := os.MkdirAll(path, 0666)

	if err != nil {
		return false
	}

	return true
}

func IsValidPath(path string) bool {
	return fs.ValidPath(path)
}