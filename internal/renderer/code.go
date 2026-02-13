package renderer

import (
	"bytes"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/v2"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// CodeExtensions maps file extensions to Chroma language identifiers.
var CodeExtensions = map[string]string{
	".ts": "typescript", ".tsx": "typescript", ".js": "javascript", ".jsx": "javascript",
	".css": "css", ".html": "html", ".json": "json", ".py": "python", ".go": "go",
	".rs": "rust", ".sh": "bash", ".yaml": "yaml", ".yml": "yaml", ".toml": "toml",
	".sql": "sql", ".java": "java", ".kt": "kotlin", ".swift": "swift", ".rb": "ruby",
	".php": "php", ".vue": "html", ".svelte": "html", ".cjs": "javascript", ".mjs": "javascript",
	".env": "bash",
}

// IsCodeFile returns true if the file should be rendered as code (any non-.md file).
func IsCodeFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext != ".md" && ext != ""
}

// RenderCodeFile renders a code file with syntax highlighting and line numbers.
func RenderCodeFile(filePath string) (*RenderResult, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	lang := CodeExtensions[ext]
	if lang == "" {
		lang = "plaintext"
	}
	basename := filepath.Base(filePath)
	lines := strings.Split(string(content), "\n")

	// Build table rows with line-by-line highlighting
	var rows strings.Builder
	for i, line := range lines {
		highlighted := highlightLine(line, lang)
		fmt.Fprintf(&rows, `<tr><td class="line-num" data-line="%d">%d</td><td class="line-content">%s</td></tr>`,
			i+1, i+1, highlighted)
		rows.WriteByte('\n')
	}

	htmlContent := fmt.Sprintf(`<div class="code-viewer">
  <div class="code-viewer-header">%s &middot; %s &middot; %d lines</div>
  <pre><code class="chroma language-%s"><table class="code-table">
%s</table></code></pre>
</div>`,
		html.EscapeString(basename), lang, len(lines), lang, rows.String())

	return &RenderResult{
		HTML:        htmlContent,
		TOC:         nil,
		Frontmatter: map[string]interface{}{},
		Title:       basename,
	}, nil
}

func highlightLine(line, lang string) string {
	if line == "" {
		line = " "
	}

	lexer := lexers.Get(lang)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get("github")
	if style == nil {
		style = styles.Fallback
	}

	formatter := chromahtml.New(
		chromahtml.WithClasses(true),
		chromahtml.PreventSurroundingPre(true),
	)

	iterator, err := lexer.Tokenise(nil, line)
	if err != nil {
		return html.EscapeString(line)
	}

	var buf bytes.Buffer
	if err := formatter.Format(&buf, style, iterator); err != nil {
		return html.EscapeString(line)
	}

	return buf.String()
}
