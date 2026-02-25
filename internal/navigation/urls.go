package navigation

import (
	"net/url"
	"path/filepath"
	"strings"
)

// ViewURL returns a clean URL for the given absolute file path relative to rootDir.
func ViewURL(absPath, rootDir string) string {
	rel, err := filepath.Rel(rootDir, absPath)
	if err != nil {
		rel = absPath
	}
	return "/" + EncodePath(rel)
}

// BrowseURL returns a clean URL for the given absolute dir path relative to rootDir.
// Returns "/" for the root directory itself.
func BrowseURL(absPath, rootDir string) string {
	rel, err := filepath.Rel(rootDir, absPath)
	if err != nil || rel == "." {
		return "/"
	}
	return "/" + EncodePath(rel)
}

// EncodePath URL-encodes each segment of a filepath while preserving slashes.
func EncodePath(p string) string {
	parts := strings.Split(filepath.ToSlash(p), "/")
	for i, part := range parts {
		parts[i] = url.PathEscape(part)
	}
	return strings.Join(parts, "/")
}
