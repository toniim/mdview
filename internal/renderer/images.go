package renderer

import (
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	inlineImgRegex = regexp.MustCompile(`!\[([^\]]*)\]\(([^)\s]+)(?:\s+"[^"]*")?\)`)
	refDefRegex    = regexp.MustCompile(`(?m)^\[([^\]]+)\]:\s*(\S+)(?:\s+"[^"]*")?$`)
)

func resolveImageSrc(src, basePath string) string {
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") || strings.HasPrefix(src, "/file") {
		return src
	}
	absPath, err := filepath.Abs(filepath.Join(basePath, src))
	if err != nil {
		return src
	}
	return "/file/" + url.PathEscape(absPath)
}

func resolveImages(markdown, basePath string) string {
	// Handle inline images: ![alt](src)
	result := inlineImgRegex.ReplaceAllStringFunc(markdown, func(match string) string {
		sub := inlineImgRegex.FindStringSubmatch(match)
		if len(sub) < 3 {
			return match
		}
		alt, src := sub[1], sub[2]
		resolved := resolveImageSrc(src, basePath)
		return "![" + alt + "](" + resolved + ")"
	})

	// Handle reference-style image definitions: [label]: src
	result = refDefRegex.ReplaceAllStringFunc(result, func(match string) string {
		sub := refDefRegex.FindStringSubmatch(match)
		if len(sub) < 3 {
			return match
		}
		label, src := sub[1], sub[2]
		resolved := resolveImageSrc(src, basePath)
		return "[" + label + "]: " + resolved
	})

	return result
}
