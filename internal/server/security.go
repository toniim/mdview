package server

import (
	"path/filepath"
	"strings"
)

func isPathSafe(filePath string, allowedDirs []string) bool {
	resolved, err := filepath.Abs(filePath)
	if err != nil {
		return false
	}

	if strings.Contains(resolved, "..") || strings.ContainsRune(filePath, 0) {
		return false
	}

	if len(allowedDirs) == 0 {
		return true
	}

	for _, dir := range allowedDirs {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			continue
		}
		if strings.HasPrefix(resolved, absDir) {
			return true
		}
	}
	return false
}
