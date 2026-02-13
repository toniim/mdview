package renderer

import (
	"html"
	"regexp"
)

var mermaidBlockRegex = regexp.MustCompile("(?s)```mermaid\n(.*?)```")

// preprocessMermaid converts ```mermaid code blocks to raw HTML <pre class="mermaid">
// before goldmark processing. This avoids conflicting with the highlighting extension's
// FencedCodeBlock renderer. goldmark passes raw HTML through with WithUnsafe().
func preprocessMermaid(markdown string) string {
	return mermaidBlockRegex.ReplaceAllStringFunc(markdown, func(match string) string {
		sub := mermaidBlockRegex.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		code := sub[1]
		return "<pre class=\"mermaid\">\n" + html.EscapeString(code) + "</pre>"
	})
}
