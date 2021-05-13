package utils

import (
	"os"
	"runtime"
	"strings"
)

func GetExecutableRootPath() string {
	os := GetSystemPlatform()

	if IsWindows(os) {
		return getWindowsExecutablePath()
	}

	return getLinuxExecutablePath()
}

func getWindowsExecutablePath() string {
	path, _ := os.Executable()
	pos := strings.LastIndex(path, "\\")

	return path[0:pos]
}

func getLinuxExecutablePath() string {
	path, _ := os.Executable()
	pos := strings.LastIndex(path, "/")

	return path[0:pos]
}

func GetSystemPlatform() string {
	return runtime.GOOS
}

func IsWindows(os string) bool {
	return os == "windows"
}

func IsLinux(os string) bool {
	return os == "linux"
}

func parseInt32(i interface{}) int32 {
	obj, ok := i.(int32)

	if !ok {
		return 0
	}

	return obj
}

func parseInt64(i interface{}) int64 {
	obj, ok := i.(int64)

	if !ok {
		return 0
	}

	return obj
}

func parseUInt32(i interface{}) uint32 {
	obj, ok := i.(uint32)

	if !ok {
		return 0
	}

	return obj
}

func parseUInt64(i interface{}) uint64 {
	obj, ok := i.(uint64)

	if !ok {
		return 0
	}

	return obj
}

func parseString(i interface{}) string {
	obj, ok := i.(string)

	if !ok {
		return ""
	}

	return obj
}