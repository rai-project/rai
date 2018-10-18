package cmd

import (
	"path/filepath"
	"runtime"
	"strings"
)

// Gets rid of volume drive label in Windows
func sanitize(name string) string {
	if len(name) > 1 && name[1] == ':' && runtime.GOOS == "windows" {
		name = name[2:]
	}

	name = filepath.Clean(name)
	name = filepath.ToSlash(name)
	for strings.HasPrefix(name, "../") {
		name = name[3:]
	}

	return name
}
