package tools

import "strings"

func ConvertToWindowsPath(path string) string {
	return strings.Replace(
		path, "/", "\\", -1)
}
