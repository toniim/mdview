package templates

import "embed"

// Assets holds the embedded filesystem containing all static assets.
// It is initialized by main.go via SetAssets before any HTTP requests.
var Assets embed.FS

// SetAssets sets the embedded filesystem. Must be called before serving.
func SetAssets(fs embed.FS) {
	Assets = fs
}
