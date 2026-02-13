package main

import (
	"embed"
	"os"

	"github.com/bilabl/mdview/cmd"
	"github.com/bilabl/mdview/internal/templates"
)

//go:embed all:assets
var assets embed.FS

func main() {
	templates.SetAssets(assets)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
