package cmd

import (
	"path/filepath"
	"runtime"
	"strings"
)

// Gets rid of volume drive label in Windows
func sanitize(name string) string {
	if runtime.GOOS == "linux" {

		name = filepath.Clean(name)
		name = filepath.ToSlash(name)
		for strings.HasPrefix(name, "../") {
			name = name[3:]
		}
	}

	return name
}
