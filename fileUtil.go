package log4z

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func getAbsUrl(relaPath string) string {
	if strings.Index(relaPath, "/") == 0 {
		return relaPath
	} else {
		currentPath := getCurrentDirectory()
		return filepath.Join(currentPath, relaPath)
	}
}
