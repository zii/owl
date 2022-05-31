package base

import (
	"os"
	"strings"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func CleanPath(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}
