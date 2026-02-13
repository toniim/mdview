//go:build ignore

package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/styles"
)

func main() {
	generateCSS("github", "assets/chroma-light.css")
	generateCSS("github-dark", "assets/chroma-dark.css")
	fmt.Println("Generated chroma CSS files")
}

func generateCSS(styleName, outPath string) {
	style := styles.Get(styleName)
	if style == nil {
		fmt.Fprintf(os.Stderr, "Style %q not found, using fallback\n", styleName)
		style = styles.Fallback
	}

	formatter := html.New(html.WithClasses(true))
	f, err := os.Create(outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating %s: %v\n", outPath, err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Fprintf(f, "/* Chroma syntax highlighting - %s style */\n", styleName)
	if err := formatter.WriteCSS(f, style); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing CSS: %v\n", err)
		os.Exit(1)
	}
}
